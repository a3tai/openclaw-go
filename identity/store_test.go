package identity

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildDeviceIdentity_V2Payload_SignatureVerifies(t *testing.T) {
	seed := make([]byte, 32)
	for i := range seed {
		seed[i] = byte(i)
	}
	priv := ed25519.NewKeyFromSeed(seed)
	pub := priv.Public().(ed25519.PublicKey)

	pubB64 := base64.RawURLEncoding.EncodeToString(pub)
	deviceID := fmt.Sprintf("%x", sha256.Sum256(pub))

	id := &Identity{
		DeviceID:        deviceID,
		PublicKeyB64URL: pubB64,
		privateKey:      priv,
	}

	p := SigningParams{
		ClientID:   "gateway-client",
		ClientMode: "backend",
		Role:       "operator",
		Scopes:     []string{"operator.admin", "sessions.read"},
		Token:      "tok_123",
		Nonce:      "nonce_456",
	}

	proto := id.BuildDeviceIdentity(p)
	if proto.ID != deviceID {
		t.Fatalf("proto.ID = %q, want %q", proto.ID, deviceID)
	}
	if proto.PublicKey != pubB64 {
		t.Fatalf("proto.PublicKey = %q, want %q", proto.PublicKey, pubB64)
	}
	if proto.Nonce != p.Nonce {
		t.Fatalf("proto.Nonce = %q, want %q", proto.Nonce, p.Nonce)
	}
	if proto.SignedAt <= 0 {
		t.Fatalf("proto.SignedAt = %d", proto.SignedAt)
	}

	scopes := strings.Join(p.Scopes, ",")
	payload := fmt.Sprintf(
		"v2|%s|%s|%s|%s|%s|%d|%s|%s",
		deviceID,
		p.ClientID,
		p.ClientMode,
		p.Role,
		scopes,
		proto.SignedAt,
		p.Token,
		p.Nonce,
	)

	sig, err := base64.RawURLEncoding.DecodeString(proto.Signature)
	if err != nil {
		t.Fatalf("decode signature: %v", err)
	}
	if !ed25519.Verify(pub, []byte(payload), sig) {
		t.Fatalf("signature did not verify")
	}
}

func TestNewIdentityFromSeed_MatchesStoreLoadOrGenerate(t *testing.T) {
	// Generate via Store, capture the seed, reconstruct via NewIdentityFromSeed,
	// and confirm the resulting Identity signs identically.
	tmp := t.TempDir()
	s, err := NewStore(tmp)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	original, err := s.LoadOrGenerate()
	if err != nil {
		t.Fatalf("LoadOrGenerate: %v", err)
	}

	seed := original.Seed()
	if len(seed) != ed25519.SeedSize {
		t.Fatalf("Seed len = %d, want %d", len(seed), ed25519.SeedSize)
	}

	rebuilt, err := NewIdentityFromSeed(seed)
	if err != nil {
		t.Fatalf("NewIdentityFromSeed: %v", err)
	}
	if rebuilt.DeviceID != original.DeviceID {
		t.Fatalf("DeviceID = %q, want %q", rebuilt.DeviceID, original.DeviceID)
	}
	if rebuilt.PublicKeyB64URL != original.PublicKeyB64URL {
		t.Fatalf("PublicKeyB64URL = %q, want %q", rebuilt.PublicKeyB64URL, original.PublicKeyB64URL)
	}

	p := SigningParams{
		ClientID:   "gateway-client",
		ClientMode: "backend",
		Role:       "operator",
		Scopes:     []string{"operator.read"},
		Token:      "tok",
		Nonce:      "nonce",
	}
	a := original.BuildDeviceIdentity(p)
	b := rebuilt.BuildDeviceIdentity(p)
	if a.ID != b.ID || a.PublicKey != b.PublicKey {
		t.Fatalf("rebuilt identity diverges from original: %+v vs %+v", a, b)
	}
}

func TestNewIdentityFromSeed_RejectsWrongLength(t *testing.T) {
	if _, err := NewIdentityFromSeed(make([]byte, 16)); err == nil {
		t.Fatal("expected error for short seed")
	}
	if _, err := NewIdentityFromSeed(nil); err == nil {
		t.Fatal("expected error for nil seed")
	}
}

func TestIdentitySeed_ReturnsCopy(t *testing.T) {
	seed := make([]byte, ed25519.SeedSize)
	for i := range seed {
		seed[i] = byte(i)
	}
	id, err := NewIdentityFromSeed(seed)
	if err != nil {
		t.Fatalf("NewIdentityFromSeed: %v", err)
	}
	got := id.Seed()
	got[0] ^= 0xff
	again := id.Seed()
	if again[0] == got[0] {
		t.Fatalf("Seed() returned shared slice; mutation leaked back")
	}
}

func TestStoreLoadOrGenerate_MigratesTruncatedDeviceID(t *testing.T) {
	tmp := t.TempDir()
	s, err := NewStore(tmp)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	// Create a deterministic identity file that simulates the older bug:
	// deviceId was sha256(pub) truncated to 16 bytes (32 hex chars).
	seed := make([]byte, 32)
	for i := range seed {
		seed[i] = byte(255 - i)
	}
	priv := ed25519.NewKeyFromSeed(seed)
	pub := priv.Public().(ed25519.PublicKey)

	pubB64 := base64.RawURLEncoding.EncodeToString(pub)
	h := sha256.Sum256(pub)
	correctID := fmt.Sprintf("%x", h[:])
	truncatedID := fmt.Sprintf("%x", h[:16])

	kp := keypairJSON{
		DeviceID:   truncatedID,
		PublicKey:  pubB64,
		PrivateKey: base64.RawURLEncoding.EncodeToString(seed),
	}
	data, _ := json.MarshalIndent(kp, "", "  ")
	fp := filepath.Join(tmp, keypairFile)
	if err := os.WriteFile(fp, append(data, '\n'), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	id, err := s.LoadOrGenerate()
	if err != nil {
		t.Fatalf("LoadOrGenerate: %v", err)
	}
	if id.DeviceID != correctID {
		t.Fatalf("DeviceID = %q, want %q", id.DeviceID, correctID)
	}

	// Ensure the on-disk file was rewritten with the corrected ID.
	onDiskRaw, err := os.ReadFile(fp)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	var kp2 keypairJSON
	if err := json.Unmarshal(onDiskRaw, &kp2); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if kp2.DeviceID != correctID {
		t.Fatalf("stored DeviceID = %q, want %q", kp2.DeviceID, correctID)
	}
	if kp2.PublicKey != pubB64 {
		t.Fatalf("stored PublicKey mismatch")
	}
}
