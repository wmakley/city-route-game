<template>
  <div class="game-board" :style="boardStyle">
    <City v-for="city in cities" :city="city" :key="city.id" />
  </div>
</template>

<script>
import City from "./City.vue";

export default {
  components: {
    City: City,
  },

  computed: {
    boardStyle() {
      return {
        width: `${this.$store.state.board.width}px`,
        height: `${this.$store.state.board.height}px`,
      };
    },

    cities() {
      return this.$store.state.cities;
    },
  },
};
</script>

<style scoped>
.game-board {
  z-index: 10;
  display: block;
  margin-top: 10px;
  border: 2px solid rgb(139, 74, 74);
  border-radius: 10px;
  background: beige;
}
</style>
