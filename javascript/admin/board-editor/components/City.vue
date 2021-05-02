<template>
  <div class="city" :style="cityStyle">{{ city.name }}</div>
</template>

<script>
export default {
  props: ["city"],

  computed: {
    cityStyle() {
      return {
        left: this.city.pos.x + "px",
        top: this.city.pos.y + "px",
        zIndex: this.city.id * 2 + 100,
      };
    },
  },
};
</script>

<style scoped>
.city {
  display: block;
  position: relative;
  width: 100px;
  border: 2px solid blue;
  border-radius: 10px;
  background: lightgreen;
  padding: 5px;
}
</style>
