import { createApp } from 'vue';
import { createStore } from 'vuex';
import App from './board-editor/App.vue';

function toInt(maybeNumber, valueIfNaN) {
	if (typeof maybeNumber === "number") {
		if (isNaN(maybeNumber)) {
			return valueIfNaN
		} else {
			return maybeNumber
		}
	}

	const numberAsInt = parseInt(maybeNumber, 10)
	if (isNaN(numberAsInt)) {
		return valueIfNaN
	}

	return numberAsInt
}

const SelectedItemTypeCity = "City";

const store = createStore({
	state() {
		return {
			board: {
				width: 800,
				height: 500
			},

			selectedItem: null,

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
			],

			nextTempCityId: -1,
		}
	},

	getters: {
		cityIsSelected(state) {
			return state.selectedItem !== null && state.selectedItem.type === SelectedItemTypeCity;
		},

		selectedCityId(state, getters) {
			if (!getters.cityIsSelected) {
				return null;
			}

			return state.selectedItem.id;
		},

		selectedCity(state, getters) {
			const id = getters.selectedCityId;
			if (id === null) {
				return null;
			}

			for (const city of state.cities) {
				if (city.id === state.selectedItem.id) {
					return city;
				}
			}

			return null;
		}
	},

	mutations: {
		setBoardWidth(state, width) {
			const widthNumber = toInt(width, 0)
			state.board.width = widthNumber
		},

		setBoardHeight(state, height) {
			const heightNumber = toInt(height, 0)
			state.board.height = heightNumber
		},

		clearSelectedItem(state) {
			state.selectedItem = null;
		},

		setSelectedCityId(state, id) {
			const idNum = toInt(id, null)
			if (idNum === null) {
				state.selectedItem = null;
				return;
			}

			state.selectedItem = {
				type: SelectedItemTypeCity,
				id: idNum,
			};
		},

		addCity(state) {
			const id = state.nextTempCityId;
			state.nextTempCityId -= 1;

			const city = {
				id: id,
				name: "New City",
				slots: [],
				pos: {
					x: 1 + id * -2,
					y: 1 + id * -2,
				},
			};

			state.cities.push(city);
			state.selectedItem = {
				type: SelectedItemTypeCity,
				id: id,
			};
		},

		deleteCity(state, id) {
			const idInt = toInt(id, null);
			if (!idInt) {
				throw new Error("id is not an integer");
			}

			const index = state.cities.findIndex(city => city.id === idInt);
			if (index === -1) {
				throw new Error(`City with ID ${idInt} not found!`);
			}

			// console.log(`Delete city ${idInt} at index: ${index}`)

			state.cities.splice(index, 1);
		},

		setCityName(state, payload) {
			const { id, name } = payload

			if (typeof name !== "string") {
				throw new Error("city name is not a string")
			}

			const city = state.cities.find(c => c.id === id)
			if (!city) {
				throw new Error(`City ${id} not found`)
			}

			city.name = name
		},

		setCityPosX(state, payload) {
			const {id, x} = payload
			const city = state.cities.find(c => c.id === id)
			if (!city) {
				throw new Error(`City ${id} not found`)
			}

			city.pos.x = x
		},

		setCityPosY(state, payload) {
			const {id, y} = payload
			const city = state.cities.find(c => c.id === id)
			if (!city) {
				throw new Error(`City ${id} not found`)
			}

			city.pos.y = y
		}
	}
});

export default () => {
	const app = createApp(App);
	app.use(store);
	return app;
};
