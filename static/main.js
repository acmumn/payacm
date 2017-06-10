$(() => {
	let amount = 0;
	const handler = StripeCheckout.configure({
	   	key: "pk_test_Mxl12rTX2BQpuAVWg6eEElFp",
		locale: "auto",
		token: onSubmit,
	});

	$("#amount").on("input", onAmountChange);
	$("#pay").click(onPay);
	const uri = URI(window.location.href);
	const query = uri.query(true);
	if(query.amount) {
		$("#amount").val(query.amount);
		$("#amount").attr("disabled", true);
		onAmountChange();
	}
	if(query.reason) {
		$("#reason").val(query.reason);
		$("#reason").attr("disabled", true);
	}

	function addCard() {
		return $("<div>")
			.addClass("jumbotron")
			.appendTo("main")
			.hide()
			.slideDown();
	}
	function getReason() {
		return $("#reason").val();
	}
	function onAmountChange() {
		const value = $("#amount").val();
		if(value.match(/^[0-9]+(\.[0-9]{2})?$/)) {
			amount = parseInt(value.replace(".", ""));
			if(value.match(/^[0-9]+$/))
				amount *= 100;

			$("#amount-group").addClass("has-success");
			$("#amount-group").removeClass("has-error");
			$("#pay").attr("disabled", false);
		} else {
			$("#amount-group").addClass("has-error");
			$("#amount-group").removeClass("has-success");
			$("#pay").attr("disabled", true);
		}
	}
	function onPay(e) {
		handler.open({
			name: "ACM UMN",
			description: getReason() || "payacmumn",
			amount: amount,
		});
		e.preventDefault();
	}
	function onSubmit(token) {
		$("#pay-form").slideUp();
		const body = {
			amount: amount,
			email: token.email,
			reason: getReason() || "",
			token: token.id,
		};
		fetch("/", {
			body: JSON.stringify(body),
			method: "POST"
		}).then(res => {
			return res.json().then(json => {
				if(res.ok)
					return json;
				else
					throw { code: res.status + " " + res.statusText, contents: json };
			});
		}).then(res => {
			alert("Success!");
		}).catch(err => {
			addCard()
				.append($("<h2>")
					.addClass("text-danger")
					.text("An Error Occurred"))
				.append($("<p>")
					.html("Contact <a href='mailto:acm@umn.edu'>acm@umn.edu</a> with the details of this error."))
				.append($("<pre>")
					.append($("<code>").text(JSON.stringify(err, null, "\t"))));
			console.error(err);
		});
	}
});
