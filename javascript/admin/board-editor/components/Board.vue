<template>
	<div class="game-board" :style="boardStyle">
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

	computed: {
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
</style>
