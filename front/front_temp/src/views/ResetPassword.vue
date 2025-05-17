<template>
    <div class="login-page">
      <div class="login-form">
        <h2 class="form-header">Введите почту для сброса пароля</h2>
        <form @submit.prevent="handleSubmit">
          <div class="login-fields">
            <inputField2
              type="text"
              id="email"
              
              :inputClass="{'with-animation': true}"
              v-model="email"
              label="Почта*"
              :error="emailError"
              :showErrors="showErrors"
              :isEmptyError="emailError === 'Поле не может быть пустым'"
              @reset-error="resetFieldError('email', 'emailError')"
            />
          </div>
          <!-- Внутри button 
          была ещё строка :disabled="isButtonDisabled" для смены стиля 
          кнопки, но в силу непонятностей с тем, в какой момент мы эту 
          кнопку должны заблочить, строка пока была убрана-->
          <button 
            type="submit"
          >
          Продолжить 
          </button>
        </form>
      </div>
    </div>
  </template>
  
  <script setup>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import InputField2 from '@/components/InputField2.vue';

const router = useRouter();
const email = ref('');
const emailError = ref('');
const showErrors = ref(false);
const isLoading = ref(false);

const isButtonDisabled = computed(() => {
  return !email.value || emailError.value || isLoading.value;
});

const validateEmail = () => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!email.value) {
    emailError.value = 'Поле не может быть пустым';
    return false;
  } else if (!emailRegex.test(email.value)) {
    emailError.value = 'Неверный почтовый адрес';
    return false;
  }
  emailError.value = '';
  return true;
};

const resetFieldError = () => {
  emailError.value = '';
};

const handleSubmit = async () => {
  showErrors.value = true;
  if (!validateEmail()) return;

  isLoading.value = true;
  
  // Имитация запроса на сервер
  await new Promise(resolve => setTimeout(resolve, 800));
  
  // Перенаправление на страницу нового пароля
  router.push('/auth/create-password');
  
  isLoading.value = false;
};
</script>
  
  
<style scoped>
.login-page {
    display: flex;
    height: 100vh;

    justify-content: center;
    align-items: center;
}

.login-form {
    width: 326px;

    margin: auto;
    padding: 0;

    font-family: "FavoritPro";
    color: var(--form-title-color);

    border: none;
}

.form-header {
    margin: 0;
    padding: 0;
    margin-bottom: 40px;

    text-align: left;
    letter-spacing: var(--form-letter-spacing-title);

    font-size: 36px;
    font-weight: 500;
    line-height: 42.73px;
}

.login-fields {
    margin-bottom: 24px;
}

button {
    width: 100%;
    padding: 16px 24px;

    font-family: inherit;
    text-align: center;
    font-size: 16px;
    line-height: 20.8px;
    font-weight: 350;

    background-color: var(--form-button-background-color);
    color: var(--form-button-color);
    letter-spacing: var(--form-letter-spacing-button);

    border: none;
}

button:hover {
    opacity: 0.5;
}

button:active {
    opacity: 1;
}

button:focus {
    box-shadow: inset 0 0 0 2px var(--form-button-border-color);
}

button:disabled {
    background-color: var(--form-button-disabled-color);
    color: var(--form-button-disabled-color);
}

</style>