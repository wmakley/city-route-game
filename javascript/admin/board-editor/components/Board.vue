<template>
	<div
		:class="cssClasses"
		:style="boardStyle"
		@dragenter="dragenter"
		@dragover="dragover"
		@dragleave="dragleave"
		@drop="drop"
	>
		<City
			v-for="(city, index) in cities"
			:city="city"
			:key="city.id"
			:index="index"
		/>
	</div>
</template>

<script>
import City from "./City.vue";

export default {
	components: {
		City: City,
	},

	data() {
		return {
			isDraggedOver: false,
		};
	},

	methods: {
		dragenter(event) {
			if (!event.dataTransfer.types.includes("text/cityid")) {
				return;
			}
			event.preventDefault();
			event.dataTransfer.dropEffect = "move";
			this.isDraggedOver = true;
		},

		dragover(event) {
			if (!event.dataTransfer.types.includes("text/cityid")) {
				return;
			}
			event.preventDefault();
			event.dataTransfer.dropEffect = "move";
			this.isDraggedOver = true;
		},

		dragleave(event) {
			this.isDraggedOver = false;
		},

		drop(event) {
			event.preventDefault();
			this.isDraggedOver = false;
			const data = event.dataTransfer.getData("text/cityid").toString();
			const cityId = parseInt(data, 10);
			// console.log("drop", event, "cityId", cityId);

			const x = Math.round(event.layerX - 60);
			const y = Math.round(event.layerY - 20);

			// console.log("new x:", x, "y:", y);

			this.$store.dispatch("setCityPosition", {
				id: cityId,
				x: x,
				y: y,
			});
		},
	},

	computed: {
		cssClasses() {
			let classes = "game-board";
			if (this.isDraggedOver) {
				classes += " dragged-over";
			}
			return classes;
		},

		boardStyle() {
			return {
				width: `${this.$store.state.board.width}px`,
				height: `${this.$store.state.board.height}px`,
			};
		},

		cities() {
			return this.$store.state.cities;
		},
	},
};
</script>

<style scoped>
.game-board {
	z-index: 1;
	display: block;
	margin-top: 10px;
	border: 2px solid rgb(139, 74, 74);
	border-radius: 10px;
	background: beige;
	overflow: auto;
	position: relative;
	top: 0;
	left: 0;

	/* ensure can scroll past the inspector */
	margin-bottom: 200px;
}

.game-board.dragged-over {
	background: cyan;
}
</style>
