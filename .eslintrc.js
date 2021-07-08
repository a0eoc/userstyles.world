module.exports = {
    parser: '@typescript-eslint/parser',
    plugins: ['@typescript-eslint', 'unicorn'],
    extends: ['plugin:@typescript-eslint/recommended'],
    rules: {
        'array-bracket-spacing': ['error', 'never'],
        'arrow-parens': ['error', 'always'],
        'block-spacing': ['error', 'always'],
        'brace-style': 'off',
        'comma-spacing': 'off',
        'curly': 'error',
        'eol-last': ['error', 'always'],
        'eqeqeq': ['error', 'smart'],
        'indent': ['error', 4, {
            'SwitchCase': 1,
        }],
        'jsx-quotes': ['error', 'prefer-double'],
        'keyword-spacing': 'off',
        'object-curly-spacing': ['error', 'never'],
        'operator-assignment': ['error', 'always'],
        'prefer-template': 'error',
        'no-debugger': 'error',
        'no-else-return': 'error',
        'no-cond-assign': 'error',
        'no-lonely-if': 'error',
        'no-multi-spaces': 'error',
        'no-implicit-coercion': 'error',
        'no-redeclare': 'off',
        'no-useless-concat': 'error',
        'no-useless-return': 'error',
        'no-trailing-spaces': 'error',
        'no-whitespace-before-property': 'error',
        'padded-blocks': ['error', 'never'],
        'prefer-exponentiation-operator': 'error',
        'prefer-regex-literals': 'error',
        'semi': 'off',
        'space-before-function-paren': ['error', {
            anonymous: 'always',
            asyncArrow: 'always',
            named: 'never',
        }],
        'space-in-parens': ['error', 'never'],
        'space-infix-ops': 'error',
        'spaced-comment': ['error', 'always', {exceptions: ['-']}],
        'quotes': 'off',
        '@typescript-eslint/array-type': ['error', {
            default: 'array-simple',
        }],
        'yoda': ['error', 'always'],
        '@typescript-eslint/brace-style': 'error',
        '@typescript-eslint/camelcase': 'off',
        '@typescript-eslint/comma-spacing': ['error', {
            before: false,
            after: true,
        }],
        '@typescript-eslint/consistent-type-assertions': 'error',
        '@typescript-eslint/consistent-type-imports': 'error',
        '@typescript-eslint/explicit-function-return-type': 'off',
        '@typescript-eslint/explicit-member-accessibility': 'off',
        '@typescript-eslint/explicit-module-boundary-types': 'off',
        '@typescript-eslint/func-call-spacing': ['error', 'never'],
        '@typescript-eslint/keyword-spacing': ['error', {
            after: true,
            before: true,
        }],
        '@typescript-eslint/member-delimiter-style': 'error',
        '@typescript-eslint/no-duplicate-imports': 'error',
        '@typescript-eslint/no-explicit-any': 'off',
        '@typescript-eslint/no-redeclare': 'error',
        '@typescript-eslint/no-object-literal-type-assertion': 'off',
        '@typescript-eslint/no-use-before-define': 'off',
        '@typescript-eslint/semi': ['error', 'always'],
        '@typescript-eslint/quotes': ['error', 'single', {
            allowTemplateLiterals: true,
            avoidEscape: true,
        }],
        'unicorn/better-regex': 'error',
        'unicorn/catch-error-name': ['error', {
            name: 'ball',
        }],
        'unicorn/consistent-destructuring': 'error',
        'unicorn/consistent-function-scoping': 'error',
        'unicorn/escape-case': 'error',
        'unicorn/no-array-push-push': 'error',
        'unicorn/no-lonely-if': 'error',
        'unicorn/no-new-array': 'error',
        'unicorn/no-zero-fractions': 'error',
        'unicorn/number-literal-case': 'error',
        'unicorn/prefer-ternary': 'error',
    },
};
