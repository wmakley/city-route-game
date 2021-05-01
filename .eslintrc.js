module.exports = {
    "env": {
		"node": true,
        "browser": true,
        "es2021": true
    },
	"exclude": [
		"node_modules",
		"javascript/board-editor/src"
	],
    "extends": [
        "eslint:recommended"
    ],
    "parserOptions": {
        "ecmaVersion": 12,
        "sourceType": "module"
    },
    "rules": {
    }
};
