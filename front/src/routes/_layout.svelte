<script>
	import { init_location_change_event } from './utils'
	import { onMount } from 'svelte'
	import Nav from '../components/Nav.svelte'

	export let segment;
	
	onMount(()=>{
		init_location_change_event();
	})

	
</script>

<style>
	main {
		position: relative;
		max-width: 56em;
		background-color: white;
		padding: 2em;
		margin: 0 auto;
		box-sizing: border-box;
	}
</style>

<Nav segment={segment}/>

<main>
	<slot></slot>
</main>