<template>
	<div class="inspector">
		<button class="btn btn-outline-secondary btn-sm float-end" @click="close">
			x Close
		</button>
		<CityInspector v-if="cityIsSelected" :city="selectedCity" />
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import CityInspector from "./CityInspector.vue";

export default {
	components: {
		CityInspector: CityInspector,
	},

	computed: {
		...mapGetters(["cityIsSelected", "selectedCity"]),
		isVisible() {
			return this.$store.selectedItem !== null;
		},
	},

	methods: {
		close() {
			this.$store.commit("clearSelectedItem");
		},
	},
};
</script>

<style scoped>
.inspector {
	position: fixed;
	width: 100%;
	left: 0;
	bottom: 0;
	padding: 10px;
	background: white;
	border-top: 2px solid gray;
}
</style>
