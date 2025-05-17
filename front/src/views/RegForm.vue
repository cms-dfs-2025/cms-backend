<template>
  <div class="login-page">
    <div class="login-form">
      <h2 class="form-header">Регистрация</h2>
      <form @submit.prevent="handleSubmit">
        <div class="login-fields">
          <InputField2
            type="text"
            id="name"
            :inputClass="{'with-animation': true}"
            v-model="name"
            label="Имя"
            :error="nameError"
            :showErrors="showErrors"
            :isEmptyError="nameError === 'Поле не может быть пустым'"
            @reset-error="resetFieldError('name')"
          />
          
          <InputField2
            type="text"
            id="email"
            :inputClass="{'with-animation': true}"
            v-model="email"
            label="Почта"
            :error="emailError"
            :showErrors="showErrors"
            :isEmptyError="emailError === 'Поле не может быть пустым'"
            @reset-error="resetFieldError('email')"
          />
          
          <InputField2
            type="password"
            id="password"
            :inputClass="{'with-animation': true}"
            v-model="password"
            label="Пароль"
            :error="passwordError || (showErrors && passwordsMismatch && confirmPassword ? ' ' : '')"
            :showErrors="showErrors"
            :isEmptyError="passwordError === 'Поле не может быть пустым'"
            @reset-error="resetFieldError('password')"
          />
          
          <InputField2
            type="password"
            id="confirmPassword"
            :inputClass="{'with-animation': true}"
            v-model="confirmPassword"
            label="Подтвердить пароль"
            :error="confirmPasswordError || (showErrors && passwordsMismatch ? 'Пароли не совпадают' : '')"
            :showErrors="showErrors"
            :isEmptyError="confirmPasswordError === 'Поле не может быть пустым'"
            @reset-error="resetFieldError('confirmPassword')"
          />
        </div>

        <div class="agreement-block">
          <input 
            type="checkbox" 
            id="agreement" 
            v-model="agreement"
            class="agreement-checkbox"
          >
          <label for="agreement" class="agreement-label">
            Я соглашаюсь на <a href='#'>обработку персональных данных</a>
          </label>
          <div v-if="showErrors && !agreement" class="agreement-error">
            Необходимо согласие на обработку данных
          </div>
        </div>

        <button type="submit" :disabled="!agreement">
          Зарегистрироваться
        </button>
        
        <div class="login-link-block">
          Уже есть аккаунт? <router-link to="/auth/login" class="login-link">Войти</router-link>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import InputField2 from '@/components/InputField2.vue';
import { useRouter } from 'vue-router';

const router = useRouter();

const name = ref('');
const email = ref('');
const password = ref('');
const confirmPassword = ref('');
const agreement = ref(false);
const nameError = ref('');
const emailError = ref('');
const passwordError = ref('');
const confirmPasswordError = ref('');
const showErrors = ref(false);

const passwordsMismatch = computed(() => {
  return password.value !== confirmPassword.value;
});

watch([password, confirmPassword], () => {
  if (showErrors.value) {
    showErrors.value = false;
    passwordError.value = '';
    confirmPasswordError.value = '';
  }
});

const validateName = () => {
  if (!name.value) {
    nameError.value = 'Поле не может быть пустым';
    return false;
  }
  nameError.value = '';
  return true;
};

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

const validatePassword = () => {
  if (!password.value) {
    passwordError.value = 'Поле не может быть пустым';
    return false;
  }
  passwordError.value = '';
  return true;
};

const validateConfirmPassword = () => {
  if (!confirmPassword.value) {
    confirmPasswordError.value = 'Поле не может быть пустым';
    return false;
  }
  if (passwordsMismatch.value) {
    return false;
  }
  confirmPasswordError.value = '';
  return true;
};

const resetFieldError = (field) => {
  if (field === 'name') nameError.value = '';
  else if (field === 'email') emailError.value = '';
  else if (field === 'password') passwordError.value = '';
  else if (field === 'confirmPassword') confirmPasswordError.value = '';
};

const handleSubmit = async () => {
  showErrors.value = true;
  const isNameValid = validateName();
  const isEmailValid = validateEmail();
  const isPasswordValid = validatePassword();
  const isConfirmPasswordValid = validateConfirmPassword();
  const isAgreementValid = agreement.value;
  
  if (!isNameValid || !isEmailValid || !isPasswordValid || !isConfirmPasswordValid || !isAgreementValid) {
    return;
  }

  try {
    const response = await fetch('http://localhost:8080/api/signup', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        handle: email.value, // Используем email как handle (?)
        password: password.value
      }),
    });

    if (!response.ok) {
      throw new Error('Ошибка регистрации');
    }

    // Перенаправляем на страницу входа после успешной регистрации
    router.push('/auth/login');
  } catch (error) {
    console.error("Ошибка регистрации:", error);
    emailError.value = 'Ошибка регистрации. Попробуйте позже.';
  }
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

.agreement-block {
  margin-bottom: 24px;
}

.agreement-checkbox {
  margin-right: 8px;
}

.agreement-label {
  font-size: 12px;
  color: #222;
}

.agreement-error {
  color: #cb3d35;
  font-size: 12px;
  margin-top: 4px;
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
  margin-bottom: 16px;
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
/*
button:disabled {
  background-color: var(--form-button-disabled-color);
  color: var(--form-button-disabled-color);
  cursor: not-allowed;
}
*/
.login-link-block {
  text-align: center;
  color: #222;
  font-size: 12px;
}

.login-link, .agreement-label a {
  color: #222;
  text-decoration: none;
}

.login-link:hover, .agreement-label a:hover {
  text-decoration: underline;
}

.agreement-checkbox {
  -webkit-appearance: none;
  -moz-appearance: none;
  appearance: none;
  width: 16px;
  height: 16px;
  border: 1px solid #222;
  border-radius: 2px;
  outline: none;
  cursor: pointer;
  position: relative;
  vertical-align: middle;
  margin-right: 8px;
}


.agreement-checkbox:checked::after {
  content: "";
  position: absolute;
  left: 4px;
  top: 1px;
  width: 5px;
  height: 10px;
  border: solid #222;
  border-width: 0 2px 2px 0;
  transform: rotate(45deg);
}

.agreement-checkbox:focus {
  box-shadow: 0 0 0 2px rgba(34, 34, 34, 0.2);
}
</style>