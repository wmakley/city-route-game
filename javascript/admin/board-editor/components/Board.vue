<template>
  <div class="game-board" :style="boardStyle"></div>
</template>

<script>
export default {
  props: {
    board: {
      width: Number,
      height: Number,
    },
  },

  data() {
    return {
      boardStyle: {
        width: this.board.width,
        height: this.board.height,
      },
    };
  },
};
</script>

<style scoped>
.game-board {
  display: block;
  margin-top: 10px;
  border: 2px solid blue;
  border-radius: 2em;
}
</style>
