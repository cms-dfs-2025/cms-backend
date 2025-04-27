<template>
  <div id="app">
    <nav>
      <!--
      <router-link to="/signup">Регистрация</router-link> |
      -->
      <router-link to="/test">Тестовая (компонент TextEitor)</router-link> |
      <router-link to="/auth/login">Авторизация</router-link> |
      <router-link to="/auth/reset-password">Сброс пароля</router-link> |
      <router-link to="/auth/create-password">Новый пароль с API</router-link> |
      <router-link to="/signup">Регистрация</router-link>
    </nav>
    <router-view />
  </div>
</template>

<script>
/*import { RouterLink } from 'vue-router';*/


</script>

<style>
@import '@/assets/styles/global.css'
</style>