<template>
    <div class="test-view">
      <h1>Тестирование компонента InputFieldType</h1>
      
      <form @submit.prevent="handleSubmit" class="test-form">
        <!-- Текстовое поле -->
        <InputFieldType
          v-model="formData.text"
          type="text"
          id="text-field"
          label="Обычный текст"
          placeholder="Введите что-нибудь"
          :error-message="errors.text"
          @blur="validateField('text')"
        />
  
        <!-- Email -->
        <InputFieldType
          v-model="formData.email"
          type="email"
          id="email-field"
          label="Email"
          placeholder="example@mail.com"
          :error-message="errors.email"
          @blur="validateEmail"
        />
  
        <!-- Пароль -->
        <InputFieldType
          v-model="formData.password"
          type="password"
          id="password-field"
          label="Пароль"
          placeholder="Не менее 6 символов"
          :error-message="errors.password"
          @blur="validatePassword"
        />
  
        <!-- Textarea -->
        <InputFieldType
          v-model="formData.bio"
          type="textarea"
          id="bio-field"
          label="Текстовое сообщение"
          placeholder=""
          :rows="4"
          :error-message="errors.bio"
          @blur="validateField('bio')"
        />
  
        <!-- Выпадающий список -->
        <InputFieldType
          v-model="formData.theme"
          type="dropdown"
          id="theme-field"
          label="Тема"
          placeholder="Выберите тему"
          :options="themeOptions"
          :error-message="errors.theme"
          @blur="validateField('theme')"
        />
  
        <!-- Загрузка файла -->
        <InputFieldType
          v-model="formData.avatar"
          type="file"
          id="file-field"
          label="Аватар"
          file-label="Выберите файл"
          accept="image/*"
          :error-message="errors.avatar"
          @change="validateFile"
        />
  
        <div class="form-actions">
          <button type="submit" :disabled="isSubmitting">
            {{ isSubmitting ? 'Отправка...' : 'Отправить' }}
          </button>
          <button type="button" @click="resetForm">
            Сбросить
          </button>
        </div>
      </form>
  
      <div class="debug-info">
        <h3>Данные формы:</h3>
        <pre>{{ formData }}</pre>
        <h3>Ошибки:</h3>
        <pre>{{ errors }}</pre>
      </div>
    </div>
  </template>
  
  <script>
  import InputFieldType from '@/components/InputFieldTypes.vue'
  import { ref } from 'vue'
  
  export default {
    name: 'TestView',
    components: { InputFieldType },
    setup() {
      const formData = ref({
        text: '',
        email: '',
        password: '',
        bio: '',
        theme: '',
        avatar: null
      })
  
      const errors = ref({
        text: '',
        email: '',
        password: '',
        bio: '',
        theme: '',
        avatar: ''
      })
  
      const isSubmitting = ref(false)
  
      const themeOptions = [
        { value: 'ai', label: 'Искусственный интеллект' },
        { value: 'diff', label: 'Дифференциальные уравнения' },
        { value: 'stat', label: 'Мат. статистика' },
        { value: 'oth', label: 'Другое' }
      ]
  
      const validateField = (field) => {
        if (!formData.value[field]?.toString().trim()) {
          errors.value[field] = 'Поле обязательно для заполнения'
        } else {
          errors.value[field] = ''
        }
      }
  
      const validateEmail = () => {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
        if (!formData.value.email) {
          errors.value.email = 'Поле обязательно для заполнения'
        } else if (!emailRegex.test(formData.value.email)) {
          errors.value.email = 'Введите корректный email'
        } else {
          errors.value.email = ''
        }
      }
  
      const validatePassword = () => {
        if (!formData.value.password) {
          errors.value.password = 'Поле обязательно для заполнения'
        } else if (formData.value.password.length < 6) {
          errors.value.password = 'Пароль должен содержать не менее 6 символов'
        } else {
          errors.value.password = ''
        }
      }
  
      const validateFile = () => {
        if (!formData.value.avatar) {
          errors.value.avatar = 'Выберите файл'
        } else if (formData.value.avatar.size > 2 * 1024 * 1024) {
          errors.value.avatar = 'Файл слишком большой (макс. 2MB)'
        } else {
          errors.value.avatar = ''
        }
      }
  
      const handleSubmit = () => {
        validateField('text')
        validateEmail()
        validatePassword()
        validateField('bio')
        validateField('theme')
        validateFile()
  
        const hasErrors = Object.values(errors.value).some(error => error)
        
        if (!hasErrors) {
          isSubmitting.value = true
          console.log('Форма отправлена:', formData.value)
          // Здесь будет вызов API
          setTimeout(() => {
            isSubmitting.value = false
            alert('Форма успешно отправлена!')
          }, 1500)
        }
      }
  
      const resetForm = () => {
        formData.value = {
          text: '',
          email: '',
          password: '',
          bio: '',
          theme: '',
          avatar: null
        }
        errors.value = {
          text: '',
          email: '',
          password: '',
          bio: '',
          theme: '',
          avatar: ''
        }
      }
  
      return {
        formData,
        errors,
        isSubmitting,
        themeOptions,
        validateField,
        validateEmail,
        validatePassword,
        validateFile,
        handleSubmit,
        resetForm
      }
    }
  }
  </script>
  
  <style scoped>
  .test-view {
    max-width: 800px;
    margin: 0 auto;
    padding: 2rem;
    font-family: Arial, sans-serif;
  }
  
  .test-form {
    background: #f9f9f9;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }
  
  .test-form > * + * {
    margin-top: 1.5rem;
  }
  
  .form-actions {
    display: flex;
    gap: 1rem;
    margin-top: 2rem;
  }
  
  button {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 4px;
    background: #4285f4;
    color: white;
    cursor: pointer;
    font-size: 1rem;
    transition: background 0.3s;
  }
  
  button:hover {
    background: #3367d6;
  }
  
  button:disabled {
    background: #cccccc;
    cursor: not-allowed;
  }
  
  button[type="button"] {
    background: #f1f1f1;
    color: #333;
  }
  
  .debug-info {
    margin-top: 3rem;
    padding: 1rem;
    background: #f0f0f0;
    border-radius: 4px;
  }
  
  pre {
    background: white;
    padding: 1rem;
    border-radius: 4px;
    overflow-x: auto;
  }
  </style>