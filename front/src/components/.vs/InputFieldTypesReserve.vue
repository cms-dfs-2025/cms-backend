<template>
    <div class="input-field" :class="fieldModifiers">
      <!-- Стандартные поля и textarea -->
      <template v-if="isStandardInput || type === 'textarea'">
        <component
          :is="type === 'textarea' ? 'textarea' : 'input'"
          class="input-field__native"
          :type="type"
          :id="id"
          :value="modelValue"
          @input="handleInput"
          @focus="handleFocus"
          @blur="handleBlur"
          :placeholder="floatingLabel ? ' ' : placeholder"
          :disabled="disabled"
          :rows="type === 'textarea' ? rows : null"
        />
        <label v-if="floatingLabel" :for="id" class="input-field__label">
          {{ label }}
        </label>
      </template>
      
      <!-- File Input -->
      <template v-else-if="type === 'file'">
        <div class="input-field__file-wrapper">
          <input
            class="input-field__native input-field__native--file"
            type="file"
            :id="id"
            @change="handleFileChange"
            :disabled="disabled"
          />
          <label :for="id" class="input-field__file-label">
            {{ fileLabel || label }}
          </label>
          <span v-if="modelValue" class="input-field__file-name">
            {{ modelValue.name }}
          </span>
        </div>
      </template>
  
      <!-- Dropdown -->
      <template v-else-if="type === 'dropdown'">
        <select
          class="input-field__native input-field__native--dropdown"
          :id="id"
          :value="modelValue"
          @change="handleInput"
          :disabled="disabled"
        >
          <option v-if="placeholder" value="" disabled>
            {{ placeholder }}
          </option>
          <option
            v-for="option in options"
            :key="option.value"
            :value="option.value"
          >
            {{ option.label }}
          </option>
        </select>
      </template>
  
      <template>
        <div class="input-field input-field--textarea" v-if="type === 'textarea'">
          <textarea
            class="input-field__native input-field__native--textarea"
            :id="id"
            :value="modelValue"
            @input="handleInput"
            @focus="handleFocus"
            @blur="handleBlur"
            :placeholder="floatingLabel ? ' ' : placeholder"
            :disabled="disabled"
            :rows="rows"
          ></textarea>
          <label v-if="floatingLabel" :for="id" class="input-field__label">
            {{ label }}
          </label>
        </div>
      </template>
      <!-- Иконка справа -->
      <span
        v-if="iconRight"
        class="input-field__icon-right"
        @click="$emit('icon-right-click')"
      >
        <slot name="icon-right">
          <component :is="iconRight" />
        </slot>
      </span>
    </div>
  </template>
  
  <script>
  export default {
    name: 'InputFieldTypes',
    props: {
      type: {
        type: String,
        default: 'text',
        validator: function(value) {
          return [
            'text',
            'email',
            'password',
            'textarea',
            'file',
            'dropdown'
          ].includes(value);
        }
      },
      modelValue: {
        type: [String, Number, File, Object],
        default: ''
      },
      id: {
        type: String,
        required: true
      },
      label: {
        type: String,
        default: ''
      },
      placeholder: {
        type: String,
        default: ''
      },
      disabled: {
        type: Boolean,
        default: false
      },
      floatingLabel: {
        type: Boolean,
        default: true
      },
      rows: {
        type: Number,
        default: 3
      },
      options: {
        type: Array,
        default: function() {
          return [];
        }
      },
      fileLabel: {
        type: String,
        default: 'Choose file'
      },
      iconRight: {
        type: [String, Object],
        default: null
      }
    },
    computed: {
      isStandardInput: function() {
        return ['text', 'email', 'password'].includes(this.type);
      },
      fieldModifiers: function() {
        return {
          'input-field--disabled': this.disabled,
          'input-field--with-icon': !!this.iconRight,
          [`input-field--${this.type}`]: true
        };
      }
    },
    methods: {
      handleInput: function(event) {
        this.$emit('update:modelValue', event.target.value);
      },
      handleFileChange: function(event) {
        this.$emit('update:modelValue', event.target.files[0]);
      },
      handleFocus: function(event) {
        this.$emit('focus', event);
      },
      handleBlur: function(event) {
        this.$emit('blur', event);
      }
    },
    emits: [
      'update:modelValue',
      'focus',
      'blur',
      'icon-right-click'
    ]
  };
  </script>
  
  <style scoped>
  /* Подключение глобальных переменных */
  .input-field {
    --input-text-color: var(--form-text-color, #222);
    /*--input-label-color: var(--form-label-color rgb(34, 34, 34));*/
    --input-label-color: var(--form-label-color, rgba(34, 34, 34, 0.5));
    --input-border-hover: var(--form-border-hover, #222);
    --input-error-color: var(--form-error-color, #cb3d35);
    font-family: var(--font-family-base);
  }
  
  /* Базовые стили */
  .input-field {
    position: relative;
    margin-bottom: 1.5rem;
  }
  
  
  .input-field__native {
    width: 100%;
    height: 55px;
    padding: 16px 0 12px;
    font-family: inherit;
    font-weight: 350;
    font-size: 16px;
    color: var(--input-text-color);
    border: none;
    border-bottom: 1px solid var(--input-label-color);
    background: transparent;
    transition: border-color 0.3s ease;
  }
  
  /* Textarea специфика 
  .input-field__native--textarea {
    height: auto;
    min-height: 100px;
    resize: vertical;
  }*/
  
  /* Плавающий лейбл */
  .input-field__label {
    position: absolute;
    top: 50%;
    left: 0;
    transform: translateY(-50%);
    pointer-events: none;
    font-size: 16px;
    color: #222;
    opacity: 0.5;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }
  /*
  .input-field__label {
    opacity: 1;
    transition: color 150ms ease;
  }*/
    
  .input-field__native:focus + .input-field__label,
  .input-field__native:not(:placeholder-shown) + .input-field__label {
    top: 0;
    transform: translateY(0);
    font-size: 12px;
  }
  
  /* File input стили */
  .input-field__file-wrapper {
    position: relative;
  }
  
  .input-field__native--file {
    opacity: 0;
    position: absolute;
    width: 0.1px;
    height: 0.1px;
  }
  
  .input-field__file-label {
    display: block;
    padding: 18px 0 16px;
    border-bottom: 1px solid var(--input-label-color);
    cursor: pointer;
    transition: border-color 0.3s ease;
  }
  
  .input-field__file-name {
    display: block;
    margin-top: 4px;
    font-size: 12px;
    color: var(--input-label-color);
  }
  
  /* Иконка справа */
  .input-field__icon-right {
    position: absolute;
    right: 0;
    top: 50%;
    transform: translateY(-50%);
    cursor: pointer;
    color: var(--input-label-color);
  }
  
  /* Специфичные стили для textarea */
  .input-field--textarea .input-field__label {
    top: 18px; /* Совпадает с padding-top textarea */
    left: 0;
    transform: none; /* Убираем центровку */
    align-items: flex-start; /* Выравнивание по верху */
  }
  
  /* Анимация при фокусе/заполнении */
  .input-field--textarea .input-field__native:focus + .input-field__label,
  .input-field--textarea .input-field__native:not(:placeholder-shown) + .input-field__label {
    top: 0;
    transform: translateY(0);
    font-size: 12px;
  }
  
  /* Дополнительные стили для textarea */
  .input-field__native--textarea {
    padding-top: 38px; /* Место для лейбла */
    min-height: 120px;
  }
  </style>