<template>
	<div v-if="city !== null" class="city-inspector">
		<p>City Inspector:</p>
		<input
			type="text"
			class="form-control"
			placeholder="City Name"
			v-model="cityName"
		/>
		<input
			type="number"
			class="form-control"
			placeholder="X"
			step="1"
			v-model.number="cityPosX"
		/>
		<input
			type="number"
			class="form-control"
			placeholder="Y"
			step="1"
			v-model.number="cityPosY"
		/>
	</div>
</template>

<script>
export default {
	props: ["city"],

	computed: {
		cityName: {
			get() {
				return this.city.name;
			},
			set(value) {
				this.$store.commit({
					type: "setCityName",
					id: this.city.id,
					name: value,
				});
			},
		},
		cityPosX: {
			get() {
				return this.city.pos.x;
			},
			set(val) {
				this.$store.commit({
					type: "setCityPosX",
					id: this.city.id,
					x: val,
				});
			},
		},
		cityPosY: {
			get() {
				return this.city.pos.y;
			},
			set(val) {
				this.$store.commit({
					type: "setCityPosY",
					id: this.city.id,
					y: val,
				});
			},
		},
	},
};
</script>
