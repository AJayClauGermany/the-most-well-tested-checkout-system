//! A micro-benchmark with Criterion: `cargo bench`.
//!
//! The checkout is tiny, so one benchmark is enough to show the tool. Criterion runs the
//! function many times, reports the time per call with a confidence interval, and compares
//! with the previous run. This is not a performance project (see docs/testing-matrix.md);
//! CI only checks that the benchmark compiles and runs (`cargo bench -- --quick`).
//!
//! `harness = false` in Cargo.toml: Criterion brings its own `main` (the macro below).

use checkout::{Basket, BasketItem, Catalog, checkout};
use criterion::{Criterion, criterion_group, criterion_main};
use std::hint::black_box;
use std::path::Path;

/// The standard catalog from spec/catalog.json.
fn load_catalog() -> Catalog {
    let path = Path::new(env!("CARGO_MANIFEST_DIR")).join("../spec/catalog.json");
    let text = std::fs::read_to_string(&path).expect("spec/catalog.json should exist");
    serde_json::from_str(&text).expect("spec/catalog.json should be a catalog")
}

fn item(sku: &str, quantity: i64) -> BasketItem {
    BasketItem {
        sku: sku.to_string(),
        quantity,
    }
}

fn one_apple() -> Basket {
    Basket {
        items: vec![item("APPLE", 1)],
        coupon: None,
        shipping: "standard".to_string(),
    }
}

/// Every product, both promotions, a coupon and express shipping.
fn everything() -> Basket {
    Basket {
        items: vec![
            item("HEADPHONES", 1),
            item("COFFEE", 6),
            item("NOTEBOOK", 12),
            item("PEN", 5),
            item("BREAD", 1),
            item("APPLE", 3),
        ],
        coupon: Some("WELCOME5".to_string()),
        shipping: "express".to_string(),
    }
}

fn checkout_benchmark(c: &mut Criterion) {
    let catalog = load_catalog();
    let small = one_apple();
    let big = everything();

    let mut group = c.benchmark_group("checkout");
    // `black_box` stops the compiler from noticing the result is unused and skipping the work.
    group.bench_function("one apple", |b| {
        b.iter(|| checkout(black_box(&small), &catalog))
    });
    group.bench_function(
        "every product, both promotions, a coupon and express shipping",
        |b| b.iter(|| checkout(black_box(&big), &catalog)),
    );
    group.finish();
}

criterion_group!(benches, checkout_benchmark);
criterion_main!(benches);
