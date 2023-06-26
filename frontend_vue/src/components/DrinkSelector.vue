<script setup lang="ts">

function test() {
    console.log('test');
}

</script>

<template>
    <div class="drink-selection">
        <svg @click="test"
             class="drink-selection__navigation-up"
             width="10%"
             height="10%"
             viewBox="0 0 50 30"
             xmlns="http://www.w3.org/2000/svg">
            <g>
                <rect x="0" y="20" width="10" height="10"/>
                <rect x="10" y="10" width="10" height="10"/>
                <rect x="20" y="0" width="10" height="10"/>
                <rect x="30" y="10" width="10" height="10"/>
                <rect x="40" y="20" width="10" height="10"/>
            </g>
        </svg>
        <span class="drink-selection__drink" @click="$emit('start')">drink</span>
        <svg @click="test"
             class="drink-selection__navigation-down"
             width="10%"
             height="10%"
             viewBox="0 0 50 30"
             xmlns="http://www.w3.org/2000/svg">
            <g>
                <rect x="0" y="20" width="10" height="10"/>
                <rect x="10" y="10" width="10" height="10"/>
                <rect x="20" y="0" width="10" height="10"/>
                <rect x="30" y="10" width="10" height="10"/>
                <rect x="40" y="20" width="10" height="10"/>
            </g>
        </svg>
    </div>
</template>

<style scoped>
.drink-selection {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    height: 100%;
    align-items: center;
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
    color: var(--barduino-dark-brown);
}
</style>