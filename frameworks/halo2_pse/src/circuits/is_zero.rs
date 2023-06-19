use halo2_proofs::{arithmetic::Field, circuit::*, plonk::*, poly::Rotation};

use crate::circuits::utils::f_values::get_one;

use super::utils::f_values::get_zero;

#[derive(Clone, Debug)]

pub struct IsZeroConfig<F> {
    pub value_inv: Column<Advice>, // value invert = 1/value
    pub is_zero_expr: Expression<F>, // if value = 0, then is_zero_expr = 1, else is_zero_expr = 0
    // We can use this is_zero_expr as a selector to trigger certain actions for example!
}

impl<F: Field> IsZeroConfig<F> {
    pub fn expr(&self) -> Expression<F> {
        self.is_zero_expr.clone()
    }
}

pub struct IsZeroChip<F: Field> {
    config: IsZeroConfig<F>,
}

impl<F: Field> IsZeroChip<F> {
    pub fn construct(config: IsZeroConfig<F>) -> Self {
        IsZeroChip { config }
    }

    // q_enable is a selector to enable the gate. q_enable is a closure
    // value is the value to be checked. Value is a closure
    pub fn configure(
        meta: &mut ConstraintSystem<F>,
        q_enable: impl FnOnce(&mut VirtualCells<'_, F>) -> Expression<F>,
        value: impl FnOnce(&mut VirtualCells<'_, F>) -> Expression<F>,
        value_inv: Column<Advice>,
    ) -> IsZeroConfig<F> {
        let mut is_zero_expr = Expression::Constant(get_zero());

        meta.create_gate("is_zero", |meta| {
            //
            // valid | value |  value_inv |  1 - value * value_inv | value * (1 - value* value_inv)
            // ------+-------+------------+------------------------+-------------------------------
            //  yes  |   x   |    1/x     |         0              |  0
            //  no   |   x   |    0       |         1              |  x
            //  yes  |   0   |    0       |         1              |  0
            //  yes  |   0   |    y       |         1              |  0

            // let's first get the value expression here from the lambda function
            let value = value(meta);
            let q_enable = q_enable(meta);
            // query value_inv from the advise colums
            let value_inv = meta.query_advice(value_inv, Rotation::cur());

            // This is the expression assignement for is_zero_expr
            is_zero_expr = Expression::Constant(get_one()) - value.clone() * value_inv;

            // there's a problem here. For example if we have a value x and a malicious prover add 0 to value_inv
            // then the prover can make the is_zero_expr = 1 - x * 0 = 1 - 0 = 1 which shouldn't be valid!
            // So we need to add a constraint to avoid that
            vec![q_enable * value * is_zero_expr.clone()]
        });

        IsZeroConfig {
            value_inv,
            is_zero_expr,
        }
    }

    // The assignment function takes the actual value, generate the inverse of that and assign it to the advice column
    pub fn assign(
        &self,
        region: &mut Region<'_, F>,
        offset: usize,
        value: Value<F>,
    ) -> Result<(), Error> {
        let value_inv = value.map(|value| value.invert().unwrap_or(get_zero()));
        region.assign_advice(|| "value inv", self.config.value_inv, offset, || value_inv)?;
        Ok(())
    }
}