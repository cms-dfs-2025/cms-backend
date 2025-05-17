<template>
    <div class="input-field">
      <input
        :type="type"
        :id="id"
        :class="{'input-error': showErrors && error}"
        :value="modelValue"
        @input="$emit('update:modelValue', $event.target.value)"
        @focus="resetError"
        :placeholder="placeholder"
      />
      <div v-show="showErrors && error" class="error-message">{{ error }}</div>
    </div>
  </template>
  
  <script>
  export default {
    props: {
      id: {
        type: String,
        required: true,
      },
      type: {
        type: String,
        default: 'text',
      },
      placeholder: {
        type: String,
        default: '',
      },
      modelValue: {
        type: String,
        default: '',
      },
      error: {
        type: String,
        default: '',
      },
      showErrors: {
        type: Boolean,
        default: false,
      },
    },
    methods: {
      resetError() {
        this.$emit('reset-error');
      },
    },
  };
  </script>
  
  <style scoped>
  input {
    width: 100%;
    margin-bottom: 8px;
    padding: 18px 0 16px;
    box-sizing: border-box;
    font-family: inherit;
    font-weight: 350;
    font-size: 16px;
    color: var(--input-text-color);
    letter-spacing: var(--form-letter-spacing-input);
    line-height: 20.8px;
    border: none;
    border-bottom: 1px solid var(--input-placeholder-color);
  }
  
  input.input-error {
    color: var(--input-warning-color);
    border-bottom-color: var(--input-warning-color);
  }
  
  input.input-error::placeholder {
    color: var(--input-warning-color);
  }
  
  .error-message {
    margin: 0;
    padding: 0;
    text-align: left;
    font-size: 12px;
    font-weight: 350;
    color: var(--input-warning-color);
  }
  </style>
