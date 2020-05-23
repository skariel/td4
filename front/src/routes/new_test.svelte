<!-- TODO: save state, clear it when submitted, add a "clear" button -->
<script>
	import { post, getUser } from './utils';
	import { onMount } from 'svelte';
	import { goto } from '@sapper/app';

	let user = null;
	let tname = "";
	let tdescr = "";
	let tcode = "";

    onMount(()=>{
        user = getUser();
    })

	async function create_test() {
		const res = await post(user, 'create_test', {
				title : tname,
				descr : tdescr,
				code  : tcode,
			})
		if (res.status == 200) {
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

<title>New test</title>

<h4>Title</h4>
<input bind:value={tname}>

<h4 style="margin-top:20px;">Description</h4>
<textarea style="width:100%; height:200px" bind:value={tdescr} />

<h4 style="margin-top:20px;">Code</h4>
<textarea style="width:100%; height:200px" bind:value={tcode} />

<div>
	{#if validate(tname, tdescr, tcode).length!=0}
		<h5 style="margin-top:20px; color:red;">{validate(tname, tdescr, tcode)}</h5> 
	{/if}
	<button on:click={()=>create_test()} disabled={validate(tname, tdescr, tcode).length != 0}>
		add test
	</button>
</div>
