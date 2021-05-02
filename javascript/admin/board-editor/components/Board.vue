<template>
  <div class="game-board" :style="boardStyle"></div>
</template>

<script>
export default {
  computed: {
    boardStyle() {
      return {
        width: `${this.$store.state.board.width}px`,
        height: `${this.$store.state.board.height}px`,
      };
    },
  },
};
</script>

<style scoped>
.game-board {
  display: block;
  margin-top: 10px;
  border: 2px solid rgb(139, 74, 74);
  border-radius: 10px;
  background: beige;
}
</style>
