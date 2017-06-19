$(() => {
	let amount = 0;

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

			if(amount < 1000) {
				$("#amount-error").show().text("payacm can only handle trans" +
					"actions of at least $10.00. Give an officer cash in per" +
					"son for any lesser amount.");
				$("#amount-group").addClass("has-error");
				$("#amount-group").removeClass("has-success");
				$("#pay").attr("disabled", true);
			} else {
				$("#amount-error").hide();
				$("#amount-group").addClass("has-success");
				$("#amount-group").removeClass("has-error");
				$("#pay").attr("disabled", false);
			}
		} else {
			$("#amount-error").show().text("Invalid amount. Please use eithe" +
				"r an integer number of dollars, or DOLLARS.CENTS notation.");
			$("#amount-group").addClass("has-error");
			$("#amount-group").removeClass("has-success");
			$("#pay").attr("disabled", true);
		}
	}
	function onPay(e) {
		fetch("/stripeKey.json").then(res => {
			return res.json();
		}).then(key => {
			const handler = StripeCheckout.configure({
	   			key: key,
				locale: "auto",
				token: onSubmit,
			});
			handler.open({
				name: "ACM UMN",
				description: getReason() || "payacmumn",
				amount: amount,
			});
		}).catch(err => {
			const body = $("<p>")
				.append("Contact ")
				.append($("<a>").attr("href", "mailto:acm@umn.edu").text("acm@umn.edu"))
				.append(" with the details of this error.");
			addCard()
				.append($("<h2>")
					.addClass("text-danger")
					.text("An Error Occurred"))
				.append(body)
				.append($("<pre>")
					.append($("<code>").text(JSON.stringify(err, null, "    "))));
			console.error(err);
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
			headers: new Headers({
				"Content-Type": "application/json",
			}),
			method: "POST"
		}).then(res => {
			return res.text().then(body => {
				if(res.ok) {
					return body;
				} else {
					throw {
						contents: body,
						status: res.status,
					};
				}
			});
		}).then(res => {
			addCard()
				.append($("<h2>").text("Succeeded"))
				.append($("<p>").text("Charge succeeded. Check your email inbox for a receipt."))
				.append($("<hr>"))
				.append($("<div>").html(res));
		}).catch(err => {
			const body = $("<p>");
			if(err.status === 502) {
				body.append("The transaction did ");
				body.append($("<strong>").text("not"));
				body.append(" complete successfully; your card should not");
				body.append(" have been charged. If this is not true,");
				body.append(" contact ");
				body.append($("<a>").attr("href", "mailto:acm@umn.edu").text("acm@umn.edu"))
				body.append(".");
			} else if(err.status === 504) {
				body.append("The transaction may have completed successfully,");
				body.append(" but an email could not be sent. Send an email to ");
				body.append($("<a>").attr("href", "mailto:acm@umn.edu").text("acm@umn.edu"))
				body.append(" with the details of the transaction and the");
				body.append(" below error.");
			} else {
				body.append("Contact ")
				body.append($("<a>").attr("href", "mailto:acm@umn.edu").text("acm@umn.edu"))
				body.append(" with the details of this error.");
			}

			addCard()
				.append($("<h2>")
					.addClass("text-danger")
					.text("An Error Occurred"))
				.append(body)
				.append($("<pre>")
					.append($("<code>").text(JSON.stringify(err, null, "    "))));
			console.error(err);
		});
	}
});
