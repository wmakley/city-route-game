<template>
  <Board :board="board" />
</template>

<script>
import Board from "./components/Board.vue";

export default {
  components: {
    Board: Board,
  },

  data() {
    return {
      board: {
        width: 500,
        height: 500,
      },
    };
  },
};
</script>
