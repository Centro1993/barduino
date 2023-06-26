<script setup lang="tsk">

</script>

<template>
    <section class="overlay">
        <div class="overlay__inner">
            <span>DRINK IN PROGRESS...</span>
        </div>
    </section>
</template>

<style scoped>
.overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: var(--overlay);
    margin: 0;
    display: flex;
    justify-content: center;
    align-items: center;
}

.overlay__inner {
    color: white;
    border: .25rem solid var(--barduino-medium-brown);
    background: var(--barduino-light-brown);
    padding: 1rem 2rem;
}
</style>