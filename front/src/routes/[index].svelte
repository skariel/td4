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

<title>TesTus</title>

<AllTests offset="{index}"></AllTests>

