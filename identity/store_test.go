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
		PrivateKey:      priv,
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
