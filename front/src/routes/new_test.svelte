<script>
	import { post, getUser, start_invalidate_cache } from './utils';
	import { onMount } from 'svelte';
	import { writable } from "svelte/store";
	import { goto } from '@sapper/app';

	// TODO: fix these to work "onmount"
	let user = null;
	const tname = writable(localStorage.getItem("tname") || "");
	tname.subscribe(val => localStorage.setItem("tname", val));
	const tdescr = writable(localStorage.getItem("tdescr") || "");
	tdescr.subscribe(val => localStorage.setItem("tdescr", val));
	const tcode = writable(localStorage.getItem("tcode") || "");
	tcode.subscribe(val => localStorage.setItem("tcode", val));

    onMount(()=>{
        user = getUser();
    })

	async function create_test() {
		const res = await post(user, 'create_test', {
				title : $tname,
				descr : $tdescr,
				code  : $tcode,
			})
		if (res.status == 200) {
			tname.set("")
			tdescr.set("")
			tcode.set("")
			start_invalidate_cache()
			goto('/')
		}
	}

	function validate(tname, tdescr, tcode) {
		if (tname.length == 0) {
			return "missing title"
		}
		if (tname.length > 256) {
			return `title is too long: ${tname.length} > 256`
		}
		if (tdescr.length == 0) {
			return "description is missing"
		}
		if (tdescr.length > 2048) {
			return `description is too long: ${tname.length} > 2048`
		}
		if (tcode.length == 0) {
			return "code is missing"
		}
		if (tcode.length > 8192) {
			return `code is too long: ${tname.length} > 8192`
		}
		return ""
	}

</script>

<style>
</style>

<svelte:head>
	<title>New test</title>
</svelte:head>


<h4>Title</h4>
<input bind:value={$tname}>

<h4 style="margin-top:20px;">Description</h4>
<textarea style="width:100%; height:200px" bind:value={$tdescr} />

<h4 style="margin-top:20px;">Code</h4>
<textarea style="width:100%; height:200px" bind:value={$tcode} />

<div>
	{#if validate($tname, $tdescr, $tcode).length!=0}
		<h5 style="margin-top:20px; color:red;">{validate($tname, $tdescr, $tcode)}</h5> 
	{/if}
	<button on:click={create_test} disabled={validate($tname, $tdescr, $tcode).length != 0}>
		add test
	</button>
</div>
