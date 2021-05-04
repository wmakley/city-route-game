<template>
	<div class="inspector">
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
	},
};
</script>

<style scoped>
.inspector {
	position: fixed;
	width: 100%;
	bottom: 0;
}
</style>
