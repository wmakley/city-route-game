<template>
	<div class="board-editor">
		<BoardSettings />
		<CityManager />
		<Inspector />
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
};
</script>
