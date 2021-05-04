<template>
	<div :class="cssClasses" :style="cityStyle" @click="selectCity">
		{{ city.name }}
	</div>
</template>

<script>
export default {
	props: {
		city: {
			id: Number,
			name: String,
			pos: {
				x: Number,
				y: Number,
			},
		},
	},

	computed: {
		cityStyle() {
			return {
				left: this.city.pos.x + "px",
				top: this.city.pos.y + "px",
				zIndex: this.city.id * 2 + 100,
			};
		},

		cssClasses() {
			if (this.$store.getters.selectedCityId === this.city.id) {
				return "city selected";
			} else {
				return "city";
			}
		},
	},

	methods: {
		selectCity() {
			this.$store.commit("setSelectedCityId", this.city.id);
		},
	},
};
</script>

<style scoped>
.city {
	display: block;
	position: relative;
	width: 100px;
	border: 2px solid blue;
	border-radius: 10px;
	background: lightgreen;
	padding: 5px;
	cursor: pointer;
}

.city.selected {
	border: 2px solid red;
}
</style>
