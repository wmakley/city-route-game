import { createApp } from 'vue';
import App from './board-editor/App.vue';

export default () => {
	const app = createApp(App);
	// Additional app setup here
	return app
}
