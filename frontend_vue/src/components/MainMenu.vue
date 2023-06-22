<script setup>

import DrinkSelection from "@/components/DrinkSelection.vue";

</script>

<template>
    <main>
        <header class="main-menu__header">
            <span class="main-menu__headline">BARDUINO</span>
            <span class="main-menu__credits">CRUSI 2023 EDITION</span>
        </header>
        <section class="main-menu__content">
            <div class="main-menu__left-column">
                <drink-selection />
            </div>
            <div class="main-menu__right-column">
                <img src="../assets/MiniMugV2.gif" alt="beer animation"/>
            </div>
        </section>
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
    flex-direction: column;
    align-items: center;
    width: 100%;
}

.main-menu__right-column {
    border: .25rem solid var(--barduino-medium-brown);
    background: var(--barduino-gray);
}

.main-menu__right-column > img {
    height: auto;
    max-width: 100%;
}

</style>