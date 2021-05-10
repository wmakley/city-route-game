<template>
	<div class="row g-3 align-items-center">
		<div class="col-auto">
			<label for="boardWidth" class="col-form-label">Width</label>
		</div>
		<div class="col-auto">
			<input
				type="number"
				id="boardWidth"
				class="form-control"
				v-model.number="boardWidth"
				step="10"
				min="0"
				max="10000"
			/>
		</div>
		<div class="col-auto">
			<label for="boardHeight" class="col-form-label">Height</label>
		</div>
		<div class="col-auto">
			<input
				type="number"
				id="boardHeight"
				class="form-control"
				v-model.number="boardHeight"
				step="10"
				min="0"
				max="10000"
			/>
		</div>
	</div>
</template>

<script>
export default {
	computed: {
		boardWidth: {
			get() {
				return this.$store.state.board.width;
			},
			set(value) {
				this.$store.dispatch("setBoardWidth", value);
			},
		},

		boardHeight: {
			get() {
				return this.$store.state.board.height;
			},
			set(value) {
				this.$store.dispatch("setBoardHeight", value);
			},
		},
	},
};
</script>
