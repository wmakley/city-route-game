<template>
	<div v-if="city !== null" class="city-inspector">
		<p>City Inspector:</p>
		<div class="row g-2 mb-2">
			<div class="col-sm-2 col-xs-3">
				<div class="input-group">
					<span class="input-group-text">ID</span>
					<input
						type="text"
						class="form-control city-id"
						:value="city.id"
						readonly
						aria-label="City ID"
						size="6"
					/>
				</div>
			</div>
			<div class="col-sm-10 col-xs-9">
				<div class="input-group">
					<span class="input-group-text">Name</span>
					<input
						type="text"
						id="city-inspector-name"
						autofocus
						class="form-control"
						placeholder="City Name"
						v-model="cityName"
						@blur="persistCity"
						aria-label="City Name"
					/>
				</div>
			</div>
		</div>
		<div class="row">
			<div class="col auto">
				<div class="input-group">
					<span class="input-group-text">X</span>
					<input
						type="number"
						class="form-control"
						placeholder="X"
						step="1"
						v-model.number="cityPosX"
						aria-label="City X Position"
					/>
					<span class="input-group-text">Y</span>
					<input
						type="number"
						class="form-control"
						placeholder="Y"
						step="1"
						v-model.number="cityPosY"
						aria-label="City Y Position"
					/>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	props: ["city"],

	methods: {
		persistCity() {
			this.$store.dispatch("persistCity", this.city.id);
		},
	},

	computed: {
		cityName: {
			get() {
				return this.city.name;
			},
			set(value) {
				this.$store.commit("setCityName", {
					id: this.city.id,
					name: value,
				});
			},
		},
		cityPosX: {
			get() {
				return this.city.position.x;
			},
			set(val) {
				this.$store.dispatch("setCityPosX", {
					id: this.city.id,
					x: val,
				});
			},
		},
		cityPosY: {
			get() {
				return this.city.position.y;
			},
			set(val) {
				this.$store.dispatch("setCityPosY", {
					id: this.city.id,
					y: val,
				});
			},
		},
	},
};
</script>
