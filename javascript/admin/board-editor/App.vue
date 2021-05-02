<template>
  <div class="board-editor">
    <BoardSettings />
    <Board />
  </div>
</template>

<script>
import Board from "./components/Board.vue";
import BoardSettings from "./components/BoardSettings.vue";

export default {
  components: {
    Board: Board,
    BoardSettings: BoardSettings,
  },
};
</script>
