<template>
	<div class="city-manager position-fixed end-0">
		<h3>Cities</h3>
		<ul class="list-unstyled">
			<li
				v-for="city in cities"
				:key="city.id"
				@click="setSelectedCityId(city.id)"
			>
				{{ city.name }}
				<button
					type="button"
					class="btn btn-sm btn-outline-secondary"
					@click="deleteCity(city.id)"
				>
					Delete
				</button>
			</li>
		</ul>
		<button
			type="button"
			class="btn btn-sm btn-outline-secondary"
			@click="createCity"
		>
			+ City
		</button>
	</div>
</template>

<script>
import { defineComponent } from "vue";
import { mapState, mapGetters, mapMutations } from "vuex";

export default defineComponent({
	computed: {
		...mapState(["cities"]),
		...mapGetters(["selectedCity", "selectedCityId"]),
	},

	methods: {
		...mapMutations(["setSelectedCityId"]),

		createCity() {
			this.$store.dispatch("createCity", {
				name: "New City",
			});
		},
	},
});
</script>

<style scoped>
.city-manager {
	z-index: 50000;
	background: white;
	border-left: 2px solid grey;
	border-top: 2px solid grey;
	border-bottom: 2px solid grey;
	border-radius: 10px;
	padding: 5px;
	margin: 0 0 0 5px;
	width: 200px;
	top: 100;
}

.city-manager li {
	cursor: pointer;
}
</style>
