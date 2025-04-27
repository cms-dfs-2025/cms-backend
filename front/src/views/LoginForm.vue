<template>
  <div class="login-page">
    <div class="login-form">
      <h2 class="form-header">Войти в систему</h2>
      <form @submit.prevent="handleSubmit">
        <div class="login-fields">
          <InputField
            type="text"
            id="email"
            :class="{'input-error': showErrors && emailError}"
            v-model="email"
            placeholder="Почта*"
            :error="emailError"
            :showErrors="showErrors"
            @reset-error="resetFieldError('email', 'emailError')"
          />
          
          <InputField
              type="password"
              id="password"
              :class="{'input-error': showErrors && passwordError}"
              v-model="password"
              placeholder="Пароль*"
              :error="passwordError"
              :showErrors="showErrors"
              @reset-error="resetFieldError('password', 'passwordError')"
          />
        </div>
        <!-- Внутри button 
        была ещё строка :disabled="isButtonDisabled" для смены стиля 
        кнопки, но в силу непонятностей с тем, в какой момент мы эту 
        кнопку должны заблочить, строка пока была убрана-->
        <button 
          type="submit"
        >
        Войти 
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
  
</script>

<script>
import InputField from '@/components/InputField.vue';

export default {
components: {
  InputField,
},
data() {
  return {
    email: '',
    password: '',
    emailError: '',
    passwordError: '',
    showErrors: false,
  };
},
computed: {
  isButtonDisabled() {
    return !this.email || !this.password || this.emailError || this.passwordError;
  },
},

methods: {
  // Проверяем, является ли поле "Почта" 
  // пустым или неверно заполненным. Выводим соответствующее предупреждение
  validateEmail() {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!this.email) {
      this.emailError = 'Поле не может быть пустым'
    } else if (!emailRegex.test(this.email)) {
      this.emailError = 'Неверный почтовый адрес'
    } else {
      this.emailError = ''
    }

  },
  // Для пароля проверяем, заполнено ли поле
  validatePassword() {
    if (!this.password) {
      this.passwordError = 'Поле не может быть пустым';
    } else {
      this.passwordError = '';
    }
  },
  handleSubmit() {
  this.showErrors = true; // Показываем сообщения об ошибке

  this.validateEmail();
  this.validatePassword();
  if (this.emailError) {
    this.setFieldErrorStyle('email')

  }
  if (this.passwordError) {
    this.setFieldErrorStyle('password')
  }
  },
  resetFieldError(fieldId, errorField) {
    this[errorField] = ''
    const field = document.getElementById(fieldId);
    field.style.color = '#222222';
    field.style.borderBottomColor = '#22222277';
  },
  setFieldErrorStyle(fieldId) {
    const field = document.getElementById(fieldId);
    if (field) {
      field.style.color = '#cb3d35';
      field.style.borderBottomColor = '#cb3d35';
    }
  },
},
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