<template>
	<div
		:class="cssClasses"
		:style="cityStyle"
		draggable="true"
		@click="selectCity"
		@dragstart="dragstart"
		@dragend="dragend"
	>
		{{ city.name }}
	</div>
</template>

<script>
export default {
	props: {
		city: {
			id: Number,
			name: String,
			position: {
				x: Number,
				y: Number,
			},
		},
		index: Number,
	},

	data() {
		return {
			isDragging: false,
		};
	},

	methods: {
		dragstart(event) {
			event.dataTransfer.setData("text/cityid", this.city.id.toString());
			event.dataTransfer.effectAllowed = "move";
			// console.log("dragstart", event);
			this.isDragging = true;
		},

		dragend(event) {
			// console.log("dragend", event);
			this.isDragging = false;
		},

		selectCity() {
			this.$store.commit("setSelectedCityId", this.city.id);
		},
	},

	computed: {
		cityStyle() {
			return {
				left: this.city.position.x + "px",
				top: this.city.position.y + "px",
				zIndex: this.index * 2 + 100,
			};
		},

		cssClasses() {
			let classes = "city";

			if (this.$store.getters.selectedCityId === this.city.id) {
				classes += " selected";
			}

			if (this.isDragging) {
				classes += " dragging";
			}

			return classes;
		},
	},
};
</script>

<style scoped>
.city {
	display: block;
	position: absolute;
	width: 120px;
	height: 40px;
	border: 2px solid blue;
	border-radius: 10px;
	background: lightgreen;
	padding: 5px;
	cursor: pointer;
	text-align: center;
}

.city.selected {
	border: 2px solid red;
}

.city.dragging {
	background: red;
}
</style>
