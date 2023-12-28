<script lang="ts" setup>

import DrinkSelection from "@/components/DrinkSelector.vue";
import Overlay from "@/components/Overlay.vue";
import {ref} from "vue";

let isBrewing = ref<boolean>();

function startDrink() {
  isBrewing.value = true;
}

</script>

<template>
  <main>
    <header class="main-menu__header">
      <span class="main-menu__headline">BARDUINO</span>
      <span class="main-menu__credits">JONAS ▪ MARTIN ▪ AARON</span>
    </header>
    <section class="main-menu__content">
      <div class="main-menu__left-column">
        <drink-selection @start="startDrink"/>
      </div>
      <div class="main-menu__right-column">
        <img alt="beer animation" src="@/assets/MiniMugV2.gif"/>
      </div>
    </section>
    <overlay v-if="isBrewing" @click="isBrewing = false"/>
  </main>
</template>

<style scoped>
main {
  border: .25rem solid var(--barduino-medium-brown);
  background: var(--barduino-light-brown);
  box-shadow: .25rem .25rem 0 0 var(--barduino-dark-brown), .5rem .5rem 0 0 var(--barduino-gray);
}

main > * {
  margin: 1rem;
}

.main-menu__header {
  color: var(--barduino-medium-brown);
}

.main-menu__headline {
  font-size: 4rem;
}

.main-menu__credits {
  font-size: 1rem;
}

.main-menu__content {
  display: flex;
  justify-content: space-between;
}

header {
  display: flex;
  align-items: end;
  justify-content: space-between;
  border-bottom: .25rem solid var(--barduino-medium-brown);
}

.main-menu__left-column {
  display: flex;
  align-items: center;
  flex-direction: column;
  width: 100%;
}

.main-menu__right-column {
  border: .25rem solid var(--barduino-medium-brown);
  background: var(--barduino-gray);
}

.main-menu__right-column > img {
  max-width: 100%;
  height: auto;
}

</style>