<!-- TODO: save state, clear it when submitted, add a "clear" button -->
<script>
	import { post, getUser } from './utils';
	import { onMount } from 'svelte';
	import { writable } from "svelte/store";
	import { goto } from '@sapper/app';

	let user = null;
	const scode = writable(localStorage.getItem("scode") || "");
	scode.subscribe(val => localStorage.setItem("scode", val));
	let test_id = 0;

    onMount(()=>{
		user = getUser();
		const url = new URL(location)
		test_id = parseInt(url.searchParams.get("test_id"))
    })

	async function create_solution() {
		const res = await post(user, 'create_solution', {
				test_code_id: test_id,
				code        : $scode,
			})
		if (res.status == 200) {
			scode.set("")
			goto('/test?id='+test_id+'&page=0')
		}
	}

	function validate(scode) {
		if (scode.length == 0) {
			return "code is missing";
		}
		if (scode.length > 8192) {
			return `code is too long: ${scode.length} > 8192`;
		}
		return "";
	}

</script>

<style>
</style>

<svelte:head>
	<title>New Solution for Test {test_id}</title>
</svelte:head>

<h1>New solution for <a href={"/test?id="+test_id}>test {test_id}</a></h1>

<h4 style="margin-top:20px;">Code</h4>
<textarea style="width:100%; height:200px" bind:value={$scode} />

<div>
	{#if validate($scode).length > 0}
		<h5 style="margin-top:20px; color:red;">{validate($scode)}</h5> 
	{/if}
	<button on:click={()=>create_solution()} disabled={validate($scode).length != 0}>
		add solution
	</button>
</div>
