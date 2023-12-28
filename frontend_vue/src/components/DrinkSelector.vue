<script lang="ts" setup>

function test() {
  console.log('test');
}

</script>

<template>
  <div class="drink-selection">
    <svg class="drink-selection__navigation-up"
         height="10%"
         viewBox="0 0 50 30"
         width="10%"
         xmlns="http://www.w3.org/2000/svg"
         @click="test">
      <g>
        <rect height="10" width="10" x="0" y="20"/>
        <rect height="10" width="10" x="10" y="10"/>
        <rect height="10" width="10" x="20" y="0"/>
        <rect height="10" width="10" x="30" y="10"/>
        <rect height="10" width="10" x="40" y="20"/>
      </g>
    </svg>
    <span class="drink-selection__drink" @click="$emit('start')">drink</span>
    <svg class="drink-selection__navigation-down"
         height="10%"
         viewBox="0 0 50 30"
         width="10%"
         xmlns="http://www.w3.org/2000/svg"
         @click="test">
      <g>
        <rect height="10" width="10" x="0" y="20"/>
        <rect height="10" width="10" x="10" y="10"/>
        <rect height="10" width="10" x="20" y="0"/>
        <rect height="10" width="10" x="30" y="10"/>
        <rect height="10" width="10" x="40" y="20"/>
      </g>
    </svg>
  </div>
</template>

<style scoped>
.drink-selection {
  display: flex;
  align-items: center;
  flex-direction: column;
  justify-content: space-between;
  height: 100%;
  margin: 4rem 0;
}

.drink-selection__navigation-down, .drink-selection__navigation-up {
  cursor: pointer;
  fill: var(--barduino-dark-brown);
}

.drink-selection__navigation-down {
  transform: rotate(180deg);
}

.drink-selection__drink {
  font-size: 3rem;
  cursor: pointer;
  color: var(--barduino-dark-brown);
}
</style>