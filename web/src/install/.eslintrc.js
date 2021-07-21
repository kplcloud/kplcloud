module.exports = {
  root: true,
  env: {
    node: true
  },
  'extends': [
    'plugin:vue/essential'
  ],
  rules: {
    "vue/valid-v-on": ["error", {
      "modifiers": []
    }],
    // ? javascript rules
    // 开发模式允许使用console
    'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    // 开发环境允许使用调试 (生产模式禁用)
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    // 允许使用 async-await
    'generator-star-spacing': 'off',
    // 禁止使用 var
    'no-var': 'error',
    // 函数名括号前不需要有空格
    'space-before-function-paren': 'off',
    // 代码块中避免多余留白
    'padded-blocks': 'off',
    // 最多出现3个连续空行
    'no-multiple-empty-lines': ['error', {
      'max': 3,
      'maxBOF': 1
    }],

    // 自定义规则
    'no-eval': 0,
    'eqeqeq': 0,
    'no-unused-vars': [
      'error',
      {
        'argsIgnorePattern': 'commit'
      }
    ],
    // 自定义规则

    // ? vue rules
    // html属性必须换行
    // 没有内容的元素需要使用闭合标签
    'vue/html-self-closing': 'off',
    'no-mixed-operators': 0,
    'vue/max-attributes-per-line': [
      2,
      {
        singleline: 5,
        multiline: {
          max: 1,
          allowFirstLine: false
        }
      }
    ],
    'vue/attribute-hyphenation': 0,
    'vue/component-name-in-template-casing': 0,
    'vue/html-closing-bracket-spacing': 0,
    'vue/singleline-html-element-content-newline': 0,
    'vue/no-unused-components': 0,
    'vue/multiline-html-element-content-newline': 0,
    'vue/no-use-v-if-with-v-for': 0,
    'vue/html-closing-bracket-newline': 0,
    'vue/no-parsing-error': 0,
    'no-tabs': 0,
    quotes: [
      2,
      'single',
      {
        avoidEscape: true,
        allowTemplateLiterals: true
      }
    ],
    semi: [
      2,
      'never',
      {
        beforeStatementContinuationChars: 'never'
      }
    ],
    'no-delete-var': 2,
    'prefer-const': [
      2,
      {
        ignoreReadBeforeAssign: false
      }
    ],
    'template-curly-spacing': 'off',
    indent: 'off'
  },
  parserOptions: {
    parser: 'babel-eslint'
  },
  overrides: [
    {
      files: [
        '**/__tests__/*.{j,t}s?(x)',
        '**/tests/unit/**/*.spec.{j,t}s?(x)'
      ],
      env: {
        jest: true
      }
    }
  ]
}
