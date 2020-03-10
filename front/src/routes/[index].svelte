<script context="module">
	export function preload({ params }) {
        return {
			index  : params.index,
		}
	}
</script>

<script>
    import AllTests from '../components/AllTests.svelte'
    export let index;
</script>

<style>
</style>

<svelte:head>
	<title>TesTus</title>
</svelte:head>
<AllTests offset="{index}"></AllTests>

