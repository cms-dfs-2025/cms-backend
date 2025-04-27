<template>
    <div class="login-page">
      <div class="login-form">
        <h2 class="form-header">Войти в систему</h2>
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
            
            <InputField2
                type="password"
                id="password"
                :class="{'input-error': showErrors && passwordError}"
                :inputClass="{'with-animation': true}"
                v-model="password"
                label="Пароль*"
                :error="passwordError"
                :showErrors="showErrors"
                :isEmptyError="passwordError === 'Поле не может быть пустым'"
                @reset-error="resetFieldError('password', 'passwordError')"
            />
          </div>
          <!-- Внутри button 
          была ещё строка :disabled="isButtonDisabled" для смены стиля 
          кнопки, но в силу непонятностей с тем, в какой момент мы эту 
          кнопку должны заблочить, строка пока была убрана-->
          
          <div v-if="showPasswordError" class="password-reset-block">
            Забыли пароль? <router-link to="/auth/reset-password" class="password-reset-link">Восстановить доступ</router-link>
          </div>

          <button 
            type="submit"
          >
          Войти 
          </button>
          <br>
          <br>
          <button @click="testHttpbin">
            Тест запроса
          </button>
        </form>
      </div>
    </div>
  </template>
  
<script setup>
    import { ref, computed } from 'vue';
 /*   import axios from 'axios';*/

    import InputField2 from '@/components/InputField2.vue'
    import api from '@/api';
    import { useRouter } from 'vue-router';
    import { watch } from 'vue';
    
    const router = useRouter()

    const email = ref('');
    const password = ref('');
    const emailError = ref('');
    const passwordError = ref('');
    const showErrors = ref(false);
    const showPasswordError = ref(false);
    const isLoading = ref(false);

    watch(password, (newValue) => {
    if (newValue) {
        showPasswordError.value = false;
    }
    });
    const isButtonDisabled = computed(() => {
      return !email.value || !password.value || emailError.value || passwordError.value || isLoading.value;
    });

    // Рандомная проверка пароля (пока бэка нет)
    const checkPassword = () => {
        return Math.random() > 0.9;
    };

    const validateEmail = () => {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!email.value) {
            emailError.value = 'Поле не может быть пустым';
            return false;
        } else if (!emailRegex.test(email.value)) {
            emailError.value = 'Неверный почтовый адрес';
            return false;
        } else {
            emailError.value = '';
            return true;
        }
    };

    const validatePassword = () => {
        if (!password.value) {
            passwordError.value = 'Поле не может быть пустым';
            return false;
        } else {
            passwordError.value = '';
            return true;
        }
    };

    const resetFieldError = (field) => {
        if (field === 'email') {
            emailError.value = '';
        } else if (field === 'password') {
            passwordError.value = '';
        }
        const fieldElement = document.getElementById(field);
        if (fieldElement) {
            fieldElement.style.color = '#222222';
            fieldElement.style.borderBottomColor = '#22222277';
        } 
    };

    const setFieldErrorStyle = (fieldId) => {
        const field = document.getElementById(fieldId);
        if (field) {
            field.style.color = '#cb3d35';
            field.style.borderBottomColor = '#cb3d35';
        }
    };
    
    // Здесь нерабочий вариант с нашей API
    const handleSubmit = async () => {
        showErrors.value = true;
        const isEmailValid = validateEmail();
        const isPasswordValid = validatePassword();
        
        if (!isEmailValid || !isPasswordValid) {
            showPasswordError.value = false;
            return;
        }

        isLoading.value = true;
        
        try {
            // Кодировка в base64
            const authHeader = "Basic " + btoa(`${email.value}:${password.value}`);
            console.log("Отправленный заголовок:", authHeader);

            // запрос с Basic Auth
            const response = await api.post('/api/login', {}, {
                headers: {
                    'Authorization': authHeader
                }
            });
            
            console.log('Успешный вход:', response.data);
            router.push('/dashboard'); // Перенаправляем после входа
            
        } catch (error) {
            if (error.response?.status === 401) {
                passwordError.value = 'Неверный email или пароль';
                showPasswordError.value = true;
            } else {
                emailError.value = 'Ошибка сервера';
                console.error("Ошибка:", error);
            }
        } finally {
            isLoading.value = false;
        }
    };    

    /* 
    //Здесь рабочий вариант с публичной API
    const handleSubmit = async () => {
        showErrors.value = true;
        const isEmailValid = validateEmail();
        const isPasswordValid = validatePassword(); // Здесь всего лишь проверка на пустоту поля
        
        if (!isEmailValid || !isPasswordValid) {
            showPasswordError.value = false; // Мы не забыли пароль, просто поле не заполнили или почту криво ввели
            return;
        }

        const isPasswordCorrect = checkPassword();

        if (!isPasswordCorrect) {
            passwordError.value = 'Пароль неверный';
            showPasswordError.value = true;
            return;
        }

        showPasswordError.value = false;

        try {
            
            await api.post('/auth/login', {
                email: email.value,
                password: password.value
            });
            
            isLoading.value = true;
            const response = await api.post('/auth/login', {
            email: email.value,
            password: password.value,
            });
            
            // Обработка успешного входа
            console.log('Успешный вход:', response.data);
            // Здесь можно сохранить токен, перенаправить пользователя и т.д.
            // localStorage.setItem('token', response.data.token);
            // router.push('/dashboard');
            
        } catch (error) {
            // Обработка ошибок сервера (использовала оператор "?.")
            if (error.response?.status === 401) {
                // Ошибка авторизации
                emailError.value = 'Неверный email';
            } else if (error.response?.status === 500) {
                // Ошибка сервера
                emailError.value = 'Сервер недоступен';
            } else {
                // Другие ошибки (нет интернета, неправильный URL и т.д.)
                emailError.value = 'Ошибка соединения';
            }
        } finally {
            isLoading.value = false;
        }
            
    };*/
    const testHttpbin = async () => {
      try {
        const response = await fetch('https://httpbin.org/post', {
          method: 'POST',
          mode: 'cors',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            testEmail: email.value,
            testPassword: password.value
          }),
        });
        
        const data = await response.json();
        console.log("Ответ от httpbin:", data);
      } catch (error) {
        console.error("Ошибка запроса:", error);
      }
    }
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

.password-reset-block {
    margin: 0;
    padding: 0;

    margin-bottom: 24px;

    color: #222;
    font-size: 12px;
    letter-spacing: -0.03em;
    line-height: 130%;
}

.password-reset-link:visited {
    color:#222;
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