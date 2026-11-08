package checkout;

import static org.junit.jupiter.api.Assertions.assertEquals;

import java.util.List;
import org.junit.jupiter.api.Test;

/** Unit tests for the small bit of logic inside {@link Basket}: the default shipping method. */
class BasketTest {

    @Test
    void shippingIsStandardWhenNotGiven() {
        Basket basket = new Basket(List.of(new BasketItem("APPLE", 1)));

        assertEquals("standard", basket.shipping());
    }

    @Test
    void shippingIsKeptWhenGiven() {
        Basket basket = new Basket(List.of(new BasketItem("APPLE", 1)), null, "express");

        assertEquals("express", basket.shipping());
    }
}
