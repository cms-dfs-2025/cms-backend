<template>
    <div class="login-page">
      <div class="login-form">
        <h2 class="form-header">Придумайте пароль</h2>
        <form @submit.prevent="handleSubmit">
          <div class="login-fields">
            <InputField2
                type="password"
                id="new-password"
                :class="{'input-error': showErrors && passwordError}"
                :inputClass="{'with-animation': false}"
                v-model="newPassword"
                label="Новый пароль"
                :error="passwordError || (showErrors && passwordsMismatch && passwordRepeat ? ' ' : '')"
                :showErrors="showErrors"
                :isEmptyError="passwordError === 'Поле не может быть пустым'"
                :withAnimation="false"
                @reset-error="resetFieldError('passwordError')"
            />
            
            <InputField2
                type="password"
                id="password-repeat"
                :class="{'input-error': showErrors && passwordRepeatError}"
                :inputClass="{'with-animation': false}"
                v-model="passwordRepeat"
                label="Повторите пароль"
                :error="passwordRepeatError || (showErrors && passwordsMismatch ? 'Пароли не совпадают' : '')"
                :showErrors="showErrors"
                :isEmptyError="passwordRepeatError === 'Поле не может быть пустым'"
                :withAnimation="false"
                @reset-error="resetFieldError('passwordRepeatError')"
            />
          </div>
          <button type="submit" :disabled="isButtonDisabled">
            Продолжить
          </button>
        </form>
      </div>
    </div>
</template>

<script setup>
    import { ref, computed, watch } from 'vue';
    import InputField2 from '@/components/InputField2.vue';
    import { useRouter } from "vue-router";
    
    const router = useRouter();
    const newPassword = ref('');
    const passwordRepeat = ref('');
    const passwordError = ref('');
    const passwordRepeatError = ref('');
    const showErrors = ref(false);
    const isLoading = ref(false);

     watch([newPassword, passwordRepeat], () => {
        // Если пользователь начал вводить заново после ошибки, чтобы не включалась подсветка на каждый его тап
        if (showErrors.value) {
            showErrors.value = false;
            passwordError.value = '';
            passwordRepeatError.value = '';
        }
    });
    
    const passwordsMismatch = computed(() => {
        return newPassword.value !== passwordRepeat.value
    });
    
    const isButtonDisabled = computed(() => {
      return !newPassword.value || !passwordRepeat.value || 
             passwordError.value || passwordRepeatError.value || 
             isLoading.value;
    });

    const validatePassword = () => {
        if (!newPassword.value) {
            passwordError.value = 'Поле не может быть пустым';
            return false;
        }
        passwordError.value = '';
        return true;
    };
    // В макете не было случая несовпадения паролей, но это подразумевается, наверное
    const validatePasswordRepeat = () => {
        if (!passwordRepeat.value) {
            passwordRepeatError.value = 'Поле не может быть пустым';
            return false;
        }
        if (showErrors.value && passwordsMismatch.value) {
            return false;
        }
        passwordRepeatError.value = '';
        return true;
    };

    const resetFieldError = (fieldError) => {
        if (fieldError === 'passwordError') {
            passwordError.value = '';
        } else if (fieldError === 'passwordRepeatError') {
            passwordRepeatError.value = '';
        }
    };

    const handleSubmit = async () => {
        showErrors.value = true;
        const isPasswordValid = validatePassword();
        const isPasswordRepeatValid = validatePasswordRepeat();

        if (!isPasswordValid || !isPasswordRepeatValid  || passwordsMismatch.value) {
            return;
        }

        isLoading.value = true;
        
        try {
            // Здесь будет запрос на смену пароля
            const response = await fetch('https://httpbin.org/post', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    newPassword: newPassword.value,
                    // Нужен токен для сброса пароля?
                }),
            });
            
            const data = await response.json();
            console.log("Ответ от сервера:", data);
            
            // Перенаправляем пользователя
            //router.push('/login');
            
        } catch (error) {
            console.error("Ошибка при смене пароля:", error);

        } finally {
            router.push('/auth/login');
            isLoading.value = false;
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
    line-height: 100%;
    color: var(--form-title-color);
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
/*
button:disabled {
    background-color: var(--form-button-disabled-color);
    color: var(--form-button-disabled-color);
}
    */
</style>