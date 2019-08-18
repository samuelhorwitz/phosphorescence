<template>
    <section>
        <h2>Script Chain Editor</h2>
        <p>Script chains are collections of scripts which are run in sequence in order to produce a playlist of tracks. Various scripts may be used to handle different phases of the playlist building from reducing the available pool of tracks to choosing subsequent ones.</p><p>This is a lighter weight version of writing a full-featured script and leverages existing script functionality while reducing flexibility. It is useful for combining scripts which already do mostly what you would like, but would work better together.</p>
        <form class="wrapper">
            <section class="chainPart">
                <header>
                    <h3>Pruners</h3>
                    <p>Pruners reduce the set of tracks which are available to the seeder and builder. Pruners may be chained, each subsequent prune reducing the set further. They are optional.</p>
                </header>
                <ul class="content">
                    <li>Pruner 1</li>
                    <li>Pruner 2</li>
                    <li>blank</li>
                    <li>blank</li>
                    <li>blank</li>
                </ul>
            </section>
            <section class="chainPart">
                <header>
                    <h3>Seeder</h3>
                    <p>The seeder picks the first track, assuming the user did not specify a seed track themselves (such as through drag and drop). It is optional; the builder's seeder function will be used as a fallback.</p>
                </header>
                <div class="content">Seeder</div>
            </section>
            <section class="chainPart">
                <header>
                    <h3 class="required">Builder</h3>
                    <p>The builder script is responsible for selecting each track in the playlist after the initial track. It may also select the initial track if there is no seeder or user-supplied seed track. It is required.</p>
                </header>
                <div class="content">Builder</div>
            </section>
        </form>
    </section>
</template>

<style scoped>
    h2 {
        margin: 0px;
        margin-top: 1em;
    }

    h3 {
        margin: 0px;
        margin-top: 0.5em;
    }

    section.chainPart {
        display: flex;
        margin: 4em 0px;
    }

    header {
        flex: 1;
        border-right: 1px solid black;
        margin-right: 1em;
        padding: 0px 1em;
        text-align: right;
    }

    header p {
        text-align: justify;
        margin-top: 1em;
    }

    .content {
        flex: 1;
    }

    ul {
        list-style: none;
        margin: 0px;
        padding: 0px;
    }

    li {
        margin: 0.5em 0px;
    }

    .required::after {
        content: '*';
        margin-left: 0.5em;
    }
</style>

<script>
    export default {
        layout: 'marketplace'
    };
</script>