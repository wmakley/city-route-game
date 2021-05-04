<template>
	<div class="board-editor">
		<BoardSettings />
		<CityManager />
		<Inspector v-if="showInspector" />
		<Board />
	</div>
</template>

<script>
import Board from "./components/Board.vue";
import BoardSettings from "./components/BoardSettings.vue";
import CityManager from "./components/CityManager.vue";
import Inspector from "./components/Inspector.vue";

export default {
	components: {
		Board: Board,
		BoardSettings: BoardSettings,
		CityManager: CityManager,
		Inspector: Inspector,
	},

	computed: {
		showInspector() {
			return this.$store.state.selectedItem !== null;
		},
	},
};
</script>
