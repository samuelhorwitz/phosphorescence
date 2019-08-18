<template>
    <header class="wrapper" :class="{min: min}">
        <div class="container">
            <h1 @click="flash">
                <nuxt-link :to="to">
                    <glitch :text="text" :enable="$store.state.loading.loading" />
                </nuxt-link>
            </h1>
        </div>
    </header>
</template>

<style scoped>
    a {
        text-decoration: none;
        color: white;
    }

    .wrapper {
        border: 7px double white;
        padding: 0em;
        margin: 0.5em;
        margin-bottom: .1em;
        display: flex;
        align-items: center;
        justify-content: center;
        grid-column: 1 / 3;
    }

    .wrapper.min {
        background-color: rgb(40, 27, 61);
        width: 6.5em;
        height: 6.5em;
    }

    h1 {
        font-family: 'Varela';
        text-transform: uppercase;
        font-weight: bold;
        font-size: 5em;
        padding: 0px;
        margin: 0px;
        position: relative;
        cursor: pointer;
    }

    .container {
        position: relative;
    }

    @media only screen and (max-height: 449px) {
        .wrapper {
            margin-left: 2em;
            grid-column: 1 / 2;
            grid-row: 1 / 2;
        }
    }

    @media only screen and (max-width: 1649px) and (max-height: 449px) {
        .wrapper {
            font-size: 6px;
        }
    }

    @media only screen and (max-width: 1249px) and (max-height: 449px) {
        .wrapper {
            font-size: 5px;
        }
    }

    @media only screen and (max-width: 1049px) and (max-height: 449px) {
        .wrapper {
            font-size: 4px;
        }
    }

    @media only screen and (max-width: 849px) and (max-height: 449px) {
        .wrapper {
            font-size: 3px;
        }
    }

    @media only screen and (max-width: 749px) and (max-height: 449px) {
        .wrapper {
            display: none;
        }
    }

    @media only screen and (min-width: 1650px) and (max-height: 449px) {
        .wrapper {
            font-size: 8px;
        }
    }

    @media only screen and (min-height: 450px) and (max-width: 1099px) {
        .wrapper {
            grid-column: 1 / 2;
        }
    }
</style>

<script>
    import glitch from '~/components/glitch';

    export default {
        components: {
            glitch
        },
        props: ['min'],
        computed: {
            text() {
                if (this.min) {
                    return 'P';
                }
                return 'Phosphorescence';
            },
            to() {
                if (this.min) {
                    return '/marketplace';
                }
                return '/';
            }
        },
        methods: {
            flash() {
                this.$store.dispatch('loading/loadFlash');
            }
        }
    };
</script>
