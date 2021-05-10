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

function storeWithInitialBoard(initialBoard) {
	return createStore({
		state() {
			return {
				board: initialBoard,
				cities: [],

				selectedItem: null,
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
			fetchCitiesSuccess(state, cities) {
				state.cities = cities;
			},

			setBoardWidth(state, width) {
				const widthNumber = toInt(width, 0)
				state.board.width = widthNumber
			},

			setBoardHeight(state, height) {
				const heightNumber = toInt(height, 0)
				state.board.height = heightNumber
			},

			persistBoardDimensionsSuccess(state, json) {
				state.board.width = json.width
				state.board.height = json.height
				state.board.updatedAt = json.updatedAt
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

			createCitySuccess(state, city) {
				state.cities.push(city);
				state.selectedItem = {
					type: SelectedItemTypeCity,
					id: city.id,
				};
			},

			// deleteCity(state, id) {
			// 	const idInt = toInt(id, null);
			// 	if (!idInt) {
			// 		throw new Error("id is not an integer");
			// 	}

			// 	const index = state.cities.findIndex(city => city.id === idInt);
			// 	if (index === -1) {
			// 		throw new Error(`City with ID ${idInt} not found!`);
			// 	}

			// 	// console.log(`Delete city ${idInt} at index: ${index}`)

			// 	state.cities.splice(index, 1);
			// },

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

				city.position.x = x
			},

			setCityPosY(state, payload) {
				const {id, y} = payload
				const city = state.cities.find(c => c.id === id)
				if (!city) {
					throw new Error(`City ${id} not found`)
				}

				city.position.y = y
			}
		},

		actions: {
			setBoardWidth({commit, dispatch}, width) {
				commit("setBoardWidth", width)
				return dispatch("persistBoardDimensions")
			},

			setBoardHeight({commit, dispatch}, height) {
				commit("setBoardHeight", height)
				return dispatch("persistBoardDimensions")
			},

			async persistBoardDimensions({commit, state}) {
				const url = `/boards/${encodeURIComponent(state.board.id)}`;
				const payload = JSON.stringify({
					id: state.board.id,
					name: state.board.name,
					width: state.board.width,
					height: state.board.height,
				})

				const response = await window.fetch(url, {
					method: "PATCH",
					body: payload,
					headers: {
						"Content-Type": "application/json; charset=utf-8",
						"Accept": "application/json",
						"X-Requested-With": "XMLHttpRequest",
					}
				});
				if (!response.ok) {
					const msg = await response.text()
					throw new Error(`Error persisting board dimensions: ${msg}`);
				}

				const responseJson = await response.json()
				commit("persistBoardDimensionsSuccess", responseJson)
			},

			async fetchCities({commit, state}) {
				const url = `/boards/${encodeURIComponent(state.board.id)}/cities/`;

				const response = await window.fetch(url, {
					method: "GET",
					headers: {
						"Accept": "application/json"
					}
				});
				if (!response.ok) {
					const msg = await response.text()
					throw new Error(`Error fetching cities: ${msg}`);
				}

				const cities = await response.json();
				commit("fetchCitiesSuccess", cities);
			},

			async createCity({commit, state}) {
				const city = {
					name: "New City",
					position: {
						x: 0,
						y: 0,
					},
				};

				const url = `/boards/${encodeURIComponent(state.board.id)}/cities/`;

				const response = await window.fetch(url, {
					method: "POST",
					body: JSON.stringify(city),
					headers: {
						"Accept": "application/json"
					}
				});

				if (!response.ok) {
					const body = await response.text();
					throw new Error(`Error creating city: ${body}`);
				}

				const createdCity = await response.json();

				commit("createCitySuccess", createdCity);
			},

			deleteCity({commit}, id) {
				const idInt = toInt(id, null);
				if (!idInt) {
					throw new Error("id is not an integer");
				}

				const index = state.cities.findIndex(city => city.id === idInt);
				if (index === -1) {
					throw new Error(`City with ID ${idInt} not found!`);
				}

				// console.log(`Delete city ${idInt} at index: ${index}`)

				// state.cities.splice(index, 1);
			},
		}
	});
}

export default (board) => {
	const app = createApp(App);
	const store = storeWithInitialBoard(board);
	app.use(store);
	return app;
};
