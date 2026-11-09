//! Test helper: turns `"$12.39"` into `1239` cents – without ever touching a floating point
//! number (`spec/features/README.md` insists, and rules §1 explains why).

use checkout::Cents;

/// `"$12.39"` or `"12.39"` → `1239`. `"$5"` → `500`. `"$5.5"` → `550`.
///
/// Returns `Err` with an explanation for anything that is not a dollar amount.
pub fn dollars_to_cents(text: &str) -> Result<Cents, String> {
    let trimmed = text.trim();
    let amount = trimmed.strip_prefix('$').unwrap_or(trimmed);
    let (dollars, cents) = amount.split_once('.').unwrap_or((amount, "0"));

    let all_digits = |part: &str| !part.is_empty() && part.bytes().all(|b| b.is_ascii_digit());
    if !all_digits(dollars) || !all_digits(cents) || cents.len() > 2 {
        return Err(format!("not a dollar amount: {text:?}"));
    }

    // "5" after the point means 50 cents, not 5: pad to two digits first.
    let cents = format!("{cents:0<2}");
    let dollars: Cents = dollars
        .parse()
        .map_err(|_| format!("not a dollar amount: {text:?}"))?;
    let cents: Cents = cents
        .parse()
        .map_err(|_| format!("not a dollar amount: {text:?}"))?;
    Ok(dollars * 100 + cents)
}
