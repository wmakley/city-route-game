import { createApp } from 'vue';
import { createStore } from 'vuex';
import App from './board-editor/App.vue';

const store = createStore({
	state() {
		return {
			board: {
				width: 800,
				height: 500
			},

			cities: [
				{
					id: 1,
					name: "Beautiful City",
					slots: [],
					pos: {
						x: 10,
						y: 10,
					},
				},
				{
					id: 2,
					name: "Ugly City",
					slots: [],
					pos: {
						x: 200,
						y: 50,
					},
				}
			]
		}
	},

	mutations: {
		setBoardWidth(state, width) {
			let widthNumber = parseInt(width, 10);
			if (isNaN(widthNumber)) {
				widthNumber = 0;
			}
			state.board.width = widthNumber
		},

		setBoardHeight(state, height) {
			let heightNumber = parseInt(height, 10);
			if (isNaN(heightNumber)) {
				heightNumber = 0;
			}
			state.board.height = heightNumber
		}
	}
});

export default () => {
	const app = createApp(App);
	app.use(store);
	return app;
};
