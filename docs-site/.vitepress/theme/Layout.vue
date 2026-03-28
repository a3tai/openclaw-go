<script setup>
import DefaultTheme from 'vitepress/theme'

const { Layout } = DefaultTheme
</script>

<template>
  <Layout>
    <template #home-hero-before>
      <div class="atmosphere" aria-hidden="true">
        <div class="atmosphere-grid"></div>
        <div class="atmosphere-orb atmosphere-orb--amber"></div>
        <div class="atmosphere-orb atmosphere-orb--blue"></div>
      </div>
    </template>
  </Layout>
</template>
