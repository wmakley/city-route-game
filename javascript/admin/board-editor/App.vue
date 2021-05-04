<template>
	<div class="board-editor">
		<BoardSettings />
		<CityList />
		<Inspector v-if="showInspector" />
		<Board />
	</div>
</template>

<script>
import Board from "./components/Board.vue";
import BoardSettings from "./components/BoardSettings.vue";
import CityList from "./components/CityList.vue";
import Inspector from "./components/Inspector.vue";

export default {
	components: {
		Board: Board,
		BoardSettings: BoardSettings,
		CityList: CityList,
		Inspector: Inspector,
	},

	computed: {
		showInspector() {
			return this.$store.state.selectedItem !== null;
		},
	},
};
</script>
